import React, { useState, useEffect } from 'react';
import { useAuthStore } from '@/store';
import { userApi } from '@/services/api';
import api from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import { User, Settings as SettingsIcon, Shield, Key, Save, Bot, Eye, EyeOff, Check, AlertCircle } from 'lucide-react';

interface UserSettings {
  full_name: string;
  email: string;
}

interface AIConfig {
  provider: string;
  openai_api_key: string;
  claude_api_key: string;
  default_model: string;
  max_tokens: number;
}

interface AIProvider {
  id: string;
  name: string;
  models: string[];
  description: string;
}

const Settings: React.FC = () => {
  const { user, setAuth } = useAuthStore();
  const [activeTab, setActiveTab] = useState('profile');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  
  // 个人信息设置
  const [userSettings, setUserSettings] = useState<UserSettings>({
    full_name: user?.full_name || '',
    email: user?.email || '',
  });

  // AI配置设置
  const [aiConfig, setAIConfig] = useState<AIConfig>({
    provider: 'openai',
    openai_api_key: '',
    claude_api_key: '',
    default_model: 'gpt-4',
    max_tokens: 2048,
  });

  const [showApiKeys, setShowApiKeys] = useState({
    openai: false,
    claude: false,
  });

  const [aiProviders] = useState<AIProvider[]>([
    {
      id: 'openai',
      name: 'OpenAI',
      models: ['gpt-4', 'gpt-4-turbo', 'gpt-3.5-turbo'],
      description: '强大的通用AI模型，适用于各种开发任务'
    },
    {
      id: 'claude',
      name: 'Claude (Anthropic)',
      models: ['claude-3-opus', 'claude-3-sonnet', 'claude-3-haiku'],
      description: '专注于代码生成和分析的AI助手'
    }
  ]);

  const [connectionStatus, setConnectionStatus] = useState<{[key: string]: 'idle' | 'testing' | 'success' | 'error'}>({
    openai: 'idle',
    claude: 'idle'
  });

  useEffect(() => {
    if (user) {
      setUserSettings({
        full_name: user.full_name,
        email: user.email,
      });
    }
    loadAIConfig();
  }, [user]);

  const loadAIConfig = async () => {
    try {
      // 获取用户AI配置
      const response = await api.get('/ai/config');
      if (response.data.success) {
        const config = response.data.data;
        setAIConfig(prev => ({
          ...prev,
          provider: config.provider || 'openai',
          default_model: config.default_model || 'gpt-4',
          max_tokens: config.max_tokens || 2048,
        }));
      }
    } catch (err) {
      // 如果获取失败，使用默认配置
      console.warn('无法获取AI配置，使用默认设置');
    }
  };

  const handleSaveProfile = async () => {
    if (!userSettings.full_name.trim()) {
      setError('姓名不能为空');
      return;
    }

    try {
      setIsLoading(true);
      setError('');
      
      const updatedUser = await userApi.updateUser({
        email: userSettings.email,
        full_name: userSettings.full_name,
      });
      
      setAuth(updatedUser, localStorage.getItem('auth_token') || '');
      setMessage('个人信息保存成功');
      
      setTimeout(() => setMessage(''), 3000);
    } catch (err: any) {
      setError(err.message || '保存失败');
    } finally {
      setIsLoading(false);
    }
  };

  const testAPIConnection = async (provider: string) => {
    const apiKey = provider === 'openai' ? aiConfig.openai_api_key : aiConfig.claude_api_key;
    
    if (!apiKey.trim()) {
      setError(`请先输入${provider === 'openai' ? 'OpenAI' : 'Claude'} API密钥`);
      return;
    }

    setConnectionStatus(prev => ({ ...prev, [provider]: 'testing' }));
    setError('');

    try {
      const response = await api.post('/ai/test-connection', {
        provider: provider,
        api_key: apiKey,
        model: aiConfig.default_model
      });
      
      if (response.data.success) {
        const result = response.data.data;
        if (result.success) {
          setConnectionStatus(prev => ({ ...prev, [provider]: 'success' }));
          setMessage(`${provider === 'openai' ? 'OpenAI' : 'Claude'} 连接测试成功 (延迟: ${result.latency}ms)`);
        } else {
          setConnectionStatus(prev => ({ ...prev, [provider]: 'error' }));
          setError(`${provider === 'openai' ? 'OpenAI' : 'Claude'} 连接测试失败: ${result.message}`);
        }
      }
      
      setTimeout(() => setMessage(''), 3000);
    } catch (err: any) {
      setConnectionStatus(prev => ({ ...prev, [provider]: 'error' }));
      const errorMessage = err.response?.data?.error || err.message;
      setError(`${provider === 'openai' ? 'OpenAI' : 'Claude'} 连接测试失败: ${errorMessage}`);
    }
  };

  const handleSaveAIConfig = async () => {
    try {
      setIsLoading(true);
      setError('');

      const response = await api.put('/ai/config', {
        provider: aiConfig.provider,
        openai_api_key: aiConfig.openai_api_key,
        claude_api_key: aiConfig.claude_api_key,
        default_model: aiConfig.default_model,
        max_tokens: aiConfig.max_tokens,
      });
      
      if (response.data.success) {
        setMessage('AI配置保存成功');
        // 重新加载配置以获取最新状态
        await loadAIConfig();
      }
      
      setTimeout(() => setMessage(''), 3000);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message;
      setError(`AI配置保存失败: ${errorMessage}`);
    } finally {
      setIsLoading(false);
    }
  };

  const tabs = [
    { id: 'profile', name: '个人信息', icon: <User className="h-4 w-4" /> },
    { id: 'ai', name: 'AI配置', icon: <Bot className="h-4 w-4" /> },
    { id: 'security', name: '安全设置', icon: <Shield className="h-4 w-4" /> },
  ];

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="bg-white rounded-lg shadow">
        {/* 头部 */}
        <div className="px-6 py-4 border-b border-gray-200">
          <h1 className="text-2xl font-semibold text-gray-900">设置</h1>
          <p className="text-sm text-gray-600 mt-1">管理您的账户设置和偏好</p>
        </div>

        <div className="flex">
          {/* 侧边栏 */}
          <div className="w-64 border-r border-gray-200">
            <nav className="p-4 space-y-1">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`w-full flex items-center px-3 py-2 text-sm font-medium rounded-md ${
                    activeTab === tab.id
                      ? 'bg-blue-50 text-blue-700 border-blue-200'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  }`}
                >
                  {tab.icon}
                  <span className="ml-3">{tab.name}</span>
                </button>
              ))}
            </nav>
          </div>

          {/* 主内容区 */}
          <div className="flex-1 p-6">
            {/* 个人信息设置 */}
            {activeTab === 'profile' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">个人信息</h3>
                  
                  <div className="grid grid-cols-1 gap-6 max-w-lg">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        姓名
                      </label>
                      <Input
                        type="text"
                        value={userSettings.full_name}
                        onChange={(e) => setUserSettings(prev => ({
                          ...prev,
                          full_name: e.target.value
                        }))}
                        placeholder="请输入您的姓名"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        邮箱
                      </label>
                      <Input
                        type="email"
                        value={userSettings.email}
                        disabled
                        className="bg-gray-50"
                      />
                      <p className="text-xs text-gray-500 mt-1">邮箱地址不可修改</p>
                    </div>
                  </div>

                  <div className="mt-6">
                    <Button
                      onClick={handleSaveProfile}
                      disabled={isLoading}
                      className="flex items-center"
                    >
                      <Save className="h-4 w-4 mr-2" />
                      {isLoading ? '保存中...' : '保存设置'}
                    </Button>
                  </div>
                </div>
              </div>
            )}

            {/* AI配置设置 */}
            {activeTab === 'ai' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">AI服务配置</h3>
                  <p className="text-sm text-gray-600 mb-6">
                    配置您的AI服务提供商和API密钥，支持OpenAI和Claude
                  </p>
                  
                  {/* 默认AI提供商选择 */}
                  <div className="bg-gray-50 rounded-lg p-4 mb-6">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                      默认AI提供商
                    </label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      {aiProviders.map((provider) => (
                        <div
                          key={provider.id}
                          className={`border-2 rounded-lg p-4 cursor-pointer transition-all ${
                            aiConfig.provider === provider.id
                              ? 'border-blue-500 bg-blue-50'
                              : 'border-gray-200 hover:border-gray-300'
                          }`}
                          onClick={() => setAIConfig(prev => ({ ...prev, provider: provider.id }))}
                        >
                          <div className="flex items-center justify-between mb-2">
                            <h4 className="font-medium text-gray-900">{provider.name}</h4>
                            {aiConfig.provider === provider.id && (
                              <Check className="h-5 w-5 text-blue-600" />
                            )}
                          </div>
                          <p className="text-sm text-gray-600 mb-3">{provider.description}</p>
                          <div className="text-xs text-gray-500">
                            支持模型: {provider.models.join(', ')}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* 动态显示选中提供商的配置 */}
                  {aiConfig.provider === 'openai' && (
                    <div className="border border-gray-200 rounded-lg p-6 mb-6">
                      <div className="flex items-center justify-between mb-4">
                        <h4 className="text-lg font-medium text-gray-900">OpenAI 配置</h4>
                        <div className="flex items-center space-x-2">
                          {connectionStatus.openai === 'success' && (
                            <span className="flex items-center text-green-600 text-sm">
                              <Check className="h-4 w-4 mr-1" />
                              已连接
                            </span>
                          )}
                          {connectionStatus.openai === 'error' && (
                            <span className="flex items-center text-red-600 text-sm">
                              <AlertCircle className="h-4 w-4 mr-1" />
                              连接失败
                            </span>
                          )}
                        </div>
                      </div>
                      
                      <div className="space-y-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            API 密钥
                          </label>
                          <div className="relative">
                            <Input
                              type={showApiKeys.openai ? "text" : "password"}
                              value={aiConfig.openai_api_key}
                              onChange={(e) => setAIConfig(prev => ({
                                ...prev,
                                openai_api_key: e.target.value
                              }))}
                              placeholder="sk-..."
                              className="pr-20"
                            />
                            <div className="absolute inset-y-0 right-0 flex items-center space-x-2 pr-3">
                              <button
                                type="button"
                                onClick={() => setShowApiKeys(prev => ({ ...prev, openai: !prev.openai }))}
                                className="text-gray-400 hover:text-gray-600"
                              >
                                {showApiKeys.openai ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                              </button>
                            </div>
                          </div>
                          <p className="text-xs text-gray-500 mt-1">
                            在 <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">OpenAI平台</a> 获取您的API密钥
                          </p>
                        </div>
                        
                        <div className="flex space-x-3">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => testAPIConnection('openai')}
                            disabled={connectionStatus.openai === 'testing' || !aiConfig.openai_api_key.trim()}
                          >
                            {connectionStatus.openai === 'testing' ? '测试中...' : '测试连接'}
                          </Button>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Claude配置 */}
                  {aiConfig.provider === 'claude' && (
                    <div className="border border-gray-200 rounded-lg p-6 mb-6">
                      <div className="flex items-center justify-between mb-4">
                        <h4 className="text-lg font-medium text-gray-900">Claude 配置</h4>
                        <div className="flex items-center space-x-2">
                          {connectionStatus.claude === 'success' && (
                            <span className="flex items-center text-green-600 text-sm">
                              <Check className="h-4 w-4 mr-1" />
                              已连接
                            </span>
                          )}
                          {connectionStatus.claude === 'error' && (
                            <span className="flex items-center text-red-600 text-sm">
                              <AlertCircle className="h-4 w-4 mr-1" />
                              连接失败
                            </span>
                          )}
                        </div>
                      </div>
                      
                      <div className="space-y-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            API 密钥
                          </label>
                          <div className="relative">
                            <Input
                              type={showApiKeys.claude ? "text" : "password"}
                              value={aiConfig.claude_api_key}
                              onChange={(e) => setAIConfig(prev => ({
                                ...prev,
                                claude_api_key: e.target.value
                              }))}
                              placeholder="sk-ant-..."
                              className="pr-20"
                            />
                            <div className="absolute inset-y-0 right-0 flex items-center space-x-2 pr-3">
                              <button
                                type="button"
                                onClick={() => setShowApiKeys(prev => ({ ...prev, claude: !prev.claude }))}
                                className="text-gray-400 hover:text-gray-600"
                              >
                                {showApiKeys.claude ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                              </button>
                            </div>
                          </div>
                          <p className="text-xs text-gray-500 mt-1">
                            在 <a href="https://console.anthropic.com/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Anthropic控制台</a> 获取您的API密钥
                          </p>
                        </div>
                        
                        <div className="flex space-x-3">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => testAPIConnection('claude')}
                            disabled={connectionStatus.claude === 'testing' || !aiConfig.claude_api_key.trim()}
                          >
                            {connectionStatus.claude === 'testing' ? '测试中...' : '测试连接'}
                          </Button>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* 高级设置 */}
                  <div className="border border-gray-200 rounded-lg p-6">
                    <h4 className="text-lg font-medium text-gray-900 mb-4">高级设置</h4>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          默认模型
                        </label>
                        <select
                          value={aiConfig.default_model}
                          onChange={(e) => setAIConfig(prev => ({ ...prev, default_model: e.target.value }))}
                          className="form-input"
                        >
                          {aiProviders.find(p => p.id === aiConfig.provider)?.models.map(model => (
                            <option key={model} value={model}>{model}</option>
                          ))}
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          最大Token数
                        </label>
                        <Input
                          type="number"
                          value={aiConfig.max_tokens}
                          onChange={(e) => setAIConfig(prev => ({
                            ...prev,
                            max_tokens: parseInt(e.target.value) || 2048
                          }))}
                          min="100"
                          max="8192"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                          控制AI响应的最大长度 (100-8192)
                        </p>
                      </div>
                    </div>
                  </div>

                  <div className="mt-6">
                    <Button
                      onClick={handleSaveAIConfig}
                      disabled={isLoading}
                      className="flex items-center"
                    >
                      <Save className="h-4 w-4 mr-2" />
                      {isLoading ? '保存中...' : '保存AI配置'}
                    </Button>
                  </div>
                </div>
              </div>
            )}

            {/* 安全设置 */}
            {activeTab === 'security' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">安全设置</h3>
                  
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <div className="flex items-center">
                      <Key className="h-5 w-5 text-blue-600 mr-2" />
                      <div>
                        <h4 className="text-sm font-medium text-blue-800">密码修改</h4>
                        <p className="text-sm text-blue-700 mt-1">
                          密码修改功能正在开发中，如需修改密码请联系管理员。
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* 消息提示 */}
            {message && (
              <div className="mt-6 bg-green-50 border border-green-200 rounded-lg p-4">
                <p className="text-sm text-green-600">{message}</p>
              </div>
            )}

            {error && (
              <div className="mt-6 bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings; 