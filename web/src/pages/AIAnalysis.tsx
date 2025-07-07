import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import './AIAnalysis.css';

interface RequirementAnalysis {
  requirement_id: string;
  project_id: string;
  raw_requirement: string;
  structured_requirement: string;
  completeness_score: number;
  analysis_status: string;
  missing_info_types: string;
  created_at: string;
}

interface Question {
  question_id: string;
  requirement_id: string;
  question_text: string;
  question_category: string;
  priority_level: number;
  answer_text?: string;
  answer_status: string;
  created_at: string;
}

const AIAnalysis: React.FC = () => {
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  
  const [requirement, setRequirement] = useState('');
  const [analysis, setAnalysis] = useState<RequirementAnalysis | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'input' | 'analysis' | 'questions'>('input');

  // 加载现有分析
  useEffect(() => {
    if (projectId) {
      loadExistingAnalysis();
    }
  }, [projectId]);

  const loadExistingAnalysis = async () => {
    try {
      const response = await fetch(`/api/ai/analysis/project/${projectId}`);
      if (response.ok) {
        const analyses = await response.json();
        if (analyses.data && analyses.data.length > 0) {
          const latestAnalysis = analyses.data[0];
          setAnalysis(latestAnalysis);
          setRequirement(latestAnalysis.raw_requirement);
          loadQuestions(latestAnalysis.requirement_id);
          setActiveTab('analysis');
        }
      }
    } catch (err) {
      console.error('加载分析失败:', err);
    }
  };

  const loadQuestions = async (requirementId: string) => {
    try {
      const response = await fetch(`/api/ai/analysis/${requirementId}/questions`);
      if (response.ok) {
        const result = await response.json();
        setQuestions(result.data || []);
      }
    } catch (err) {
      console.error('加载问题失败:', err);
    }
  };

  const analyzeRequirement = async () => {
    if (!requirement.trim()) {
      setError('请输入需求描述');
      return;
    }

    if (!projectId) {
      setError('项目ID不存在');
      return;
    }

    setIsAnalyzing(true);
    setError('');

    try {
      const response = await fetch('/api/ai/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          project_id: projectId,
          requirement: requirement,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        setAnalysis(result.data);
        setActiveTab('analysis');
        
        // 稍后加载生成的问题
        setTimeout(() => {
          loadQuestions(result.data.requirement_id);
        }, 2000);
      } else {
        const errorResult = await response.json();
        setError(errorResult.error || '分析失败');
      }
    } catch (err) {
      setError('网络错误，请重试');
    } finally {
      setIsAnalyzing(false);
    }
  };

  const answerQuestion = async (questionId: string, answer: string) => {
    try {
      const response = await fetch(`/api/ai/questions/${questionId}/answer`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          question_id: questionId,
          answer: answer,
        }),
      });

      if (response.ok) {
        // 更新问题状态
        setQuestions(prev => prev.map(q => 
          q.question_id === questionId 
            ? { ...q, answer_text: answer, answer_status: 'answered' }
            : q
        ));
      } else {
        setError('回答保存失败');
      }
    } catch (err) {
      setError('网络错误，请重试');
    }
  };

  const generatePUML = async (diagramType: string) => {
    if (!analysis) return;

    try {
      const response = await fetch('/api/ai/puml/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          analysis_id: analysis.requirement_id,
          diagram_type: diagramType,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        // 导航到PUML查看页面
        navigate(`/project/${projectId}/puml/${result.data.diagram_id}`);
      } else {
        setError('PUML生成失败');
      }
    } catch (err) {
      setError('网络错误，请重试');
    }
  };

  const generateDocument = async () => {
    if (!analysis) return;

    try {
      const response = await fetch('/api/ai/document/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          analysis_id: analysis.requirement_id,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        // 导航到文档查看页面
        navigate(`/project/${projectId}/document/${result.data.document_id}`);
      } else {
        setError('文档生成失败');
      }
    } catch (err) {
      setError('网络错误，请重试');
    }
  };

  const renderStructuredRequirement = (structuredReq: string) => {
    try {
      const parsed = JSON.parse(structuredReq);
      return (
        <div className="structured-requirement">
          <h4>结构化需求</h4>
          {parsed.core_functions && (
            <div className="requirement-section">
              <h5>核心功能</h5>
              <ul>
                {parsed.core_functions.map((func: string, index: number) => (
                  <li key={index}>{func}</li>
                ))}
              </ul>
            </div>
          )}
          {parsed.roles && (
            <div className="requirement-section">
              <h5>用户角色</h5>
              <ul>
                {parsed.roles.map((role: string, index: number) => (
                  <li key={index}>{role}</li>
                ))}
              </ul>
            </div>
          )}
          {parsed.business_processes && (
            <div className="requirement-section">
              <h5>业务流程</h5>
              <ul>
                {parsed.business_processes.map((process: string, index: number) => (
                  <li key={index}>{process}</li>
                ))}
              </ul>
            </div>
          )}
          {parsed.data_entities && (
            <div className="requirement-section">
              <h5>数据实体</h5>
              <ul>
                {parsed.data_entities.map((entity: string, index: number) => (
                  <li key={index}>{entity}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      );
    } catch {
      return <div>无法解析结构化需求</div>;
    }
  };

  return (
    <div className="ai-analysis-page">
      <div className="page-header">
        <h1>AI需求分析</h1>
        <button 
          className="back-button"
          onClick={() => navigate(`/project/${projectId}`)}
        >
          返回项目
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="tabs">
        <button 
          className={`tab ${activeTab === 'input' ? 'active' : ''}`}
          onClick={() => setActiveTab('input')}
        >
          需求输入
        </button>
        <button 
          className={`tab ${activeTab === 'analysis' ? 'active' : ''}`}
          onClick={() => setActiveTab('analysis')}
          disabled={!analysis}
        >
          分析结果
        </button>
        <button 
          className={`tab ${activeTab === 'questions' ? 'active' : ''}`}
          onClick={() => setActiveTab('questions')}
          disabled={questions.length === 0}
        >
          补充问题 ({questions.filter(q => q.answer_status === 'pending').length})
        </button>
      </div>

      <div className="tab-content">
        {activeTab === 'input' && (
          <div className="input-section">
            <h3>请描述您的项目需求</h3>
            <textarea
              value={requirement}
              onChange={(e) => setRequirement(e.target.value)}
              placeholder="请详细描述您的项目需求，包括功能、目标用户、技术要求等..."
              rows={10}
              className="requirement-input"
            />
            <button 
              onClick={analyzeRequirement}
              disabled={isAnalyzing || !requirement.trim()}
              className="analyze-button"
            >
              {isAnalyzing ? '分析中...' : '开始AI分析'}
            </button>
          </div>
        )}

        {activeTab === 'analysis' && analysis && (
          <div className="analysis-section">
            <div className="analysis-header">
              <h3>分析结果</h3>
              <div className="completeness-score">
                完整度评分: {analysis.completeness_score}/100
              </div>
            </div>

            <div className="analysis-content">
              <div className="original-requirement">
                <h4>原始需求</h4>
                <p>{analysis.raw_requirement}</p>
              </div>

              {analysis.structured_requirement && 
                renderStructuredRequirement(analysis.structured_requirement)}

              <div className="action-buttons">
                <h4>生成内容</h4>
                <div className="button-group">
                  <button onClick={() => generatePUML('class')}>
                    生成类图
                  </button>
                  <button onClick={() => generatePUML('sequence')}>
                    生成时序图
                  </button>
                  <button onClick={() => generatePUML('usecase')}>
                    生成用例图
                  </button>
                  <button onClick={() => generateDocument()}>
                    生成技术文档
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'questions' && (
          <div className="questions-section">
            <h3>补充问题</h3>
            <p>AI需要以下信息来完善需求分析:</p>
            
            <div className="questions-list">
              {questions.map((question) => (
                <QuestionCard
                  key={question.question_id}
                  question={question}
                  onAnswer={answerQuestion}
                />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

interface QuestionCardProps {
  question: Question;
  onAnswer: (questionId: string, answer: string) => void;
}

const QuestionCard: React.FC<QuestionCardProps> = ({ question, onAnswer }) => {
  const [answer, setAnswer] = useState(question.answer_text || '');
  const [isEditing, setIsEditing] = useState(question.answer_status === 'pending');

  const handleSubmit = () => {
    if (answer.trim()) {
      onAnswer(question.question_id, answer);
      setIsEditing(false);
    }
  };

  const getCategoryColor = (category: string) => {
    const colors: { [key: string]: string } = {
      functionality: '#007bff',
      technical: '#28a745',
      business: '#ffc107',
      ui_ux: '#17a2b8',
      integration: '#6f42c1'
    };
    return colors[category] || '#6c757d';
  };

  return (
    <div className="question-card">
      <div className="question-header">
        <span 
          className="question-category"
          style={{ backgroundColor: getCategoryColor(question.question_category) }}
        >
          {question.question_category}
        </span>
        <span className="priority-level">
          优先级: {question.priority_level}
        </span>
      </div>
      
      <div className="question-text">
        {question.question_text}
      </div>

      {isEditing ? (
        <div className="answer-input">
          <textarea
            value={answer}
            onChange={(e) => setAnswer(e.target.value)}
            placeholder="请输入您的回答..."
            rows={3}
          />
          <div className="answer-buttons">
            <button onClick={handleSubmit}>提交答案</button>
            <button onClick={() => setIsEditing(false)}>取消</button>
          </div>
        </div>
      ) : question.answer_status === 'answered' ? (
        <div className="answer-display">
          <strong>您的回答:</strong>
          <p>{question.answer_text}</p>
          <button onClick={() => setIsEditing(true)}>修改答案</button>
        </div>
      ) : (
        <button 
          className="answer-button"
          onClick={() => setIsEditing(true)}
        >
          回答问题
        </button>
      )}
    </div>
  );
};

export default AIAnalysis; 