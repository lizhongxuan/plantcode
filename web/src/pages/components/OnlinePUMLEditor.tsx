import React, { useState, useEffect, useRef, useCallback } from 'react';
import './OnlinePUMLEditor.css';

interface OnlinePUMLEditorProps {
  value?: string;
  onChange?: (val: string) => void;
  initialCode?: string;
  onClose?: () => void;
  onSave?: (val: string) => void;
  readOnly?: boolean;
  mode?: 'edit' | 'preview' | 'split';
  style?: React.CSSProperties;
}

const OnlinePUMLEditor: React.FC<OnlinePUMLEditorProps> = ({ value, onChange, initialCode, onClose, onSave, readOnly, mode = 'split', style }) => {
  const isControlled = typeof value === 'string' && typeof onChange === 'function';
  const [pumlCode, setPumlCode] = useState(initialCode || value || '');
  const [svgContent, setSvgContent] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [errorType, setErrorType] = useState<'parse' | 'network' | ''>('');
  const [errorLines, setErrorLines] = useState<number[]>([]);
  const editorRef = useRef<HTMLTextAreaElement>(null);
  const previewRef = useRef<HTMLDivElement>(null);
  const [previewScale, setPreviewScale] = useState(1);
  const [panPosition, setPanPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const minScale = 0.1;
  const maxScale = 3.0;
  const scaleStep = 0.1;

  const handleRender = useCallback(async () => {
    setIsLoading(true);
    setError('');
    setErrorType('');
    setErrorLines([]);
    try {
      const response = await fetch('/api/puml/render-online', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: pumlCode }),
      });
      const data = await response.json();
      if (data.success && data.imageData) {
        setSvgContent(data.imageData);
      } else {
        const errorMessage = data.message || data.error || '渲染失败';
        setError('puml代码解析错误');
        setErrorType('parse');
        setSvgContent('');
        // 解析错误信息中的行号
        console.log('错误信息:', errorMessage); // 调试信息
        const lineNumbers = parseErrorLineNumbers(errorMessage);
        console.log('提取的行号:', lineNumbers); // 调试信息
        
        // 如果没有提取到行号但有错误，可能PlantUML服务没有返回具体行号
        // 这种情况下我们可以尝试通过其他方式检测错误
        if (lineNumbers.length === 0) {
          // 可以添加简单的语法检查来标识可能的错误行
          const possibleErrorLines = detectPossibleErrors(pumlCode);
          setErrorLines(possibleErrorLines);
          console.log('检测到的可能错误行:', possibleErrorLines);
        } else {
          setErrorLines(lineNumbers);
        }
      }
    } catch (err) {
      setError('连接puml服务网络出错');
      setErrorType('network');
      setSvgContent('');
      setErrorLines([]);
    } finally {
      setIsLoading(false);
    }
  }, [pumlCode]);

  // 本地检测可能的错误行
  const detectPossibleErrors = (code: string): number[] => {
    const errorLines: number[] = [];
    const lines = code.split('\n');
    
    let hasStart = false;
    let hasEnd = false;
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim();
      const lineNum = i + 1;
      
      // 检查基本语法错误
      if (line.startsWith('@startuml')) {
        if (hasStart) {
          errorLines.push(lineNum); // 重复的开始标记
        }
        hasStart = true;
      }
      
      if (line.startsWith('@enduml')) {
        if (hasEnd) {
          errorLines.push(lineNum); // 重复的结束标记
        }
        hasEnd = true;
      }
      
      // 检查括号不匹配
      const openBraces = (line.match(/\{/g) || []).length;
      const closeBraces = (line.match(/\}/g) || []).length;
      if (openBraces !== closeBraces && line.includes('{') && line.includes('}')) {
        errorLines.push(lineNum);
      }
      
      // 检查常见的语法错误
      if (line.includes('->') && !line.includes(':') && !line.includes('[') && line.length > 5) {
        // 可能缺少标签的箭头
      }
      
      // 检查无效字符或格式
      if (line.includes('  ->  ') || line.includes('-->') && !line.includes(':')) {
        errorLines.push(lineNum);
      }
    }
    
    // 检查缺少开始或结束标记
    if (!hasStart && lines.length > 1) {
      errorLines.push(1); // 第一行缺少@startuml
    }
    if (!hasEnd && lines.length > 1) {
      errorLines.push(lines.length); // 最后一行缺少@enduml
    }
    
    return [...new Set(errorLines)];
  };

  // 解析错误信息中的行号
  const parseErrorLineNumbers = (errorMessage: string): number[] => {
    const lineNumbers: number[] = [];
    // 匹配"行 数字:"的模式
    const regex = /行\s*(\d+)/g;
    let match;
    while ((match = regex.exec(errorMessage)) !== null) {
      const lineNum = parseInt(match[1], 10);
      if (!isNaN(lineNum)) {
        lineNumbers.push(lineNum);
      }
    }
    // 也匹配"line 数字"的模式（PlantUML英文错误信息）
    const englishRegex = /line\s*(\d+)/gi;
    while ((match = englishRegex.exec(errorMessage)) !== null) {
      const lineNum = parseInt(match[1], 10);
      if (!isNaN(lineNum)) {
        lineNumbers.push(lineNum);
      }
    }
    // 尝试匹配其他可能的格式
    const otherRegex = /\bat\s+line\s+(\d+)|error\s+at\s+line\s+(\d+)|\[\s*line\s+(\d+)\s*\]/gi;
    while ((match = otherRegex.exec(errorMessage)) !== null) {
      const lineNum = parseInt(match[1] || match[2] || match[3], 10);
      if (!isNaN(lineNum)) {
        lineNumbers.push(lineNum);
      }
    }
    return [...new Set(lineNumbers)]; // 去重
  };

  // 计算错误行的位置
  const calculateErrorLinePositions = useCallback(() => {
    const editor = editorRef.current;
    if (!editor || errorLines.length === 0) return [];

    const computedStyle = getComputedStyle(editor);
    const lineHeight = parseFloat(computedStyle.lineHeight) || 22;
    const fontSize = parseFloat(computedStyle.fontSize) || 15;
    const containerPaddingTop = 16; // 对应CSS中的padding

    return errorLines.map(lineNum => {
      if (lineNum > 0) {
        const top = containerPaddingTop + (lineNum - 1) * lineHeight;
        return {
          lineNumber: lineNum,
          top: top,
          height: lineHeight
        };
      }
      return null;
    }).filter(Boolean);
  }, [errorLines]);

  const errorLinePositions = calculateErrorLinePositions();

  const handleSave = () => {
    if (onSave) onSave(pumlCode);
    if (onClose) onClose();
  };
  
  // Initial render
  useEffect(() => {
    handleRender();
  }, [handleRender]);

  // Handle Ctrl+Enter or Cmd+Enter for rendering
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
        event.preventDefault();
        handleRender();
      }
    };
    const editor = editorRef.current;
    editor?.addEventListener('keydown', handleKeyDown as any);
    return () => {
      editor?.removeEventListener('keydown', handleKeyDown as any);
    };
  }, [handleRender]);

  // SVG click handling effect
  useEffect(() => {
    const jumpToLine = (lineNumber: number) => {
        const editor = editorRef.current;
        if (!editor) return;

        const lines = editor.value.split('\n');
        if (lineNumber > 0 && lineNumber <= lines.length) {
            let charPosition = 0;
            for (let i = 0; i < lineNumber - 1; i++) {
                charPosition += lines[i].length + 1; // +1 for newline
            }
            
            editor.focus();
            // 选中整行代码以显示高亮
            editor.setSelectionRange(charPosition, charPosition + lines[lineNumber - 1].length);

            // Scroll to the selection
            const lineHeight = parseFloat(getComputedStyle(editor).lineHeight) || 20;
            const scrollTop = Math.max(0, (lineNumber - 5) * lineHeight);
            editor.scrollTop = scrollTop;
        }
    };

    const jumpToCodeElement = (elementName: string) => {
        const editor = editorRef.current;
        if (!editor) return;

        const code = editor.value;
        const lines = code.split('\n');
        const searchName = elementName.toLowerCase().trim();

        if (!searchName) return;

        // Try exact match first
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i].toLowerCase().trim();
            if (line === searchName || line.includes(`"${searchName}"`) || line.includes(`'${searchName}'`)) {
                jumpToLine(i + 1);
                return;
            }
        }
        
        // Then try partial match
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i].toLowerCase();
            if (line.includes(searchName)) {
                jumpToLine(i + 1);
                return;
            }
        }
    };

    const extractElementName = (element: Element): string => {
        // Try to get text content from the element
        let name = element.textContent?.trim() || '';
        
        // If no text content, try to find text elements within
        if (!name) {
            const textEl = element.querySelector('text');
            if (textEl) {
                name = textEl.textContent?.trim() || '';
            }
        }
        
        // If still no name, try to get title or id
        if (!name) {
            name = element.getAttribute('title') || element.getAttribute('id') || '';
        }
        
        return name;
    };

    const previewEl = previewRef.current;
    if (!previewEl || !svgContent || mode === 'edit') return;
    
    // Wait for the dangerouslySetInnerHTML to render
    const timeoutId = setTimeout(() => {
      const svg = previewEl.querySelector('svg');
      if (!svg) return;

      // Select clickable elements
      const clickableElements = svg.querySelectorAll('text, rect, ellipse, circle, polygon, path, g[class*="cluster"], g[class*="node"]');
      
      const clickHandler = (event: Event) => {
          event.preventDefault();
          event.stopPropagation();
          
          const elementName = extractElementName(event.currentTarget as Element);
          console.log('Clicked element:', elementName); // Debug log
          if (elementName) {
              jumpToCodeElement(elementName);
          }
      };

      const mouseEnterHandler = (event: Event) => {
        const target = event.currentTarget as SVGElement;
        target.style.opacity = '0.7';
        target.style.cursor = 'pointer';
      };

      const mouseLeaveHandler = (event: Event) => {
        const target = event.currentTarget as SVGElement;
        target.style.opacity = '1';
      };

      // Add event listeners to all clickable elements
      clickableElements.forEach(element => {
          element.addEventListener('click', clickHandler);
          element.addEventListener('mouseenter', mouseEnterHandler);
          element.addEventListener('mouseleave', mouseLeaveHandler);
      });

      // Cleanup function
      return () => {
        clickableElements.forEach(element => {
          element.removeEventListener('click', clickHandler);
          element.removeEventListener('mouseenter', mouseEnterHandler);
          element.removeEventListener('mouseleave', mouseLeaveHandler);
        });
      };
    }, 100);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [svgContent, mode]);

  useEffect(() => {
    if (isControlled && typeof value === 'string') {
      setPumlCode(value);
    }
  }, [value, isControlled]);

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (isControlled && onChange) {
      onChange(e.target.value);
    } else {
      setPumlCode(e.target.value);
    }
  };

  const handleZoomIn = () => {
    setPreviewScale(s => Math.min(maxScale, Math.round((s + scaleStep) * 10) / 10));
  };
  const handleZoomOut = () => {
    setPreviewScale(s => Math.max(minScale, Math.round((s - scaleStep) * 10) / 10));
  };
  const handleResetZoom = () => {
    setPreviewScale(1);
    setPanPosition({ x: 0, y: 0 });
  };
  const handleFitToScreen = () => {
    const previewEl = previewRef.current;
    if (!previewEl) return;
    
    const svg = previewEl.querySelector('svg');
    if (!svg) return;
    
    const containerRect = previewEl.getBoundingClientRect();
    const svgRect = svg.getBoundingClientRect();
    
    const scaleX = (containerRect.width - 40) / svgRect.width;
    const scaleY = (containerRect.height - 40) / svgRect.height;
    const optimalScale = Math.min(scaleX, scaleY, 1);
    
    setPreviewScale(Math.max(minScale, Math.min(maxScale, optimalScale)));
    setPanPosition({ x: 0, y: 0 });
  };
  
  // 处理鼠标滚轮缩放
  const handleWheel = useCallback((e: WheelEvent) => {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault();
      const delta = e.deltaY > 0 ? -scaleStep : scaleStep;
      setPreviewScale(s => Math.max(minScale, Math.min(maxScale, Math.round((s + delta) * 10) / 10)));
    }
  }, [minScale, maxScale, scaleStep]);
  
  // 处理拖拽平移
  const handleMouseDown = (e: React.MouseEvent) => {
    if (e.button === 0) { // 左键拖拽
      setIsDragging(true);
      setDragStart({ x: e.clientX - panPosition.x, y: e.clientY - panPosition.y });
    }
  };
  
  const handleMouseMove = useCallback((e: MouseEvent) => {
    if (isDragging) {
      setPanPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y
      });
    }
  }, [isDragging, dragStart]);
  
  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
  }, []);
  
  // 添加事件监听器
  useEffect(() => {
    const previewEl = previewRef.current;
    if (!previewEl) return;
    
    previewEl.addEventListener('wheel', handleWheel, { passive: false });
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    
    return () => {
      previewEl.removeEventListener('wheel', handleWheel);
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [handleWheel, handleMouseMove, handleMouseUp]);
  
  // 重置缩放位置当切换图表时
  useEffect(() => {
    setPanPosition({ x: 0, y: 0 });
    setPreviewScale(1);
  }, [svgContent]);

  // 监听代码变化和错误行变化，重新计算位置
  useEffect(() => {
    if (errorLines.length > 0) {
      // 延迟计算，确保DOM已更新
      const timeoutId = setTimeout(() => {
        // 强制重新渲染以更新错误行位置
        const positions = calculateErrorLinePositions();
        // 这里可以添加额外的逻辑，比如自动滚动到第一个错误行
        if (positions.length > 0) {
          const editor = editorRef.current;
          if (editor) {
            const firstErrorLine = Math.min(...errorLines);
            const lineHeight = parseFloat(getComputedStyle(editor).lineHeight) || 21;
            const scrollTop = Math.max(0, (firstErrorLine - 3) * lineHeight);
            editor.scrollTop = scrollTop;
          }
        }
      }, 100);
      return () => clearTimeout(timeoutId);
    }
  }, [errorLines, pumlCode, calculateErrorLinePositions]);

  return (
    <div className={onClose ? "online-puml-editor-modal" : "online-puml-editor-embed"} style={{ height: '100%', minHeight: 0, background: '#fff', ...style }}>
      <div className="online-puml-editor-container" style={{ height: '100%', minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {onClose && (
          <header className="ope-header">
            <h1>PlantUML 在线编辑器</h1>
            <div>
              {onSave && <button onClick={handleSave} className="ope-button ope-save-btn">保存并关闭</button>}
              <button onClick={onClose} className="ope-close-btn">&times;</button>
            </div>
          </header>
        )}
        <main
          className="ope-main"
          style={{ flex: 1, minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden', position: 'relative' }}
        >
          {mode === 'edit' && (
            <div className="ope-panel ope-editor-panel" style={{ flex: 1, minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <div className="ope-panel-header">
                <h2>PUML 代码</h2>
              </div>
              <div className="ope-textarea-container" style={{ flex: 1, position: 'relative' }}>
                {/* 错误行高亮 */}
                {errorLinePositions.map((pos: any, index: number) => (
                  <React.Fragment key={`error-edit-${pos.lineNumber}-${index}`}>
                    <div
                      className="ope-error-highlight"
                      style={{
                        top: pos.top,
                        height: pos.height,
                      }}
                    />
                    <div
                      className="ope-error-underline"
                      style={{
                        top: pos.top + pos.height - 2,
                      }}
                    />
                  </React.Fragment>
                ))}
                <textarea
                  ref={editorRef}
                  value={isControlled ? value : pumlCode}
                  onChange={handleInputChange}
                  className="ope-textarea"
                  placeholder="输入 PlantUML 代码..."
                  readOnly={readOnly}
                  style={{ flex: 1, minHeight: 0, resize: 'none', background: '#fff', position: 'relative', zIndex: 3 }}
                />
              </div>
            </div>
          )}
          {mode === 'preview' && (
            <div className="ope-panel ope-preview-panel" style={{ flex: 1, minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <div className="ope-panel-header" style={{ position: 'relative' }}>
                <h2>预览</h2>
                <div style={{ position: 'absolute', right: 16, top: '50%', transform: 'translateY(-50%)', display: 'flex', alignItems: 'center', gap: 4 }}>
                  <button onClick={handleZoomOut} style={{ width: 28, height: 28, fontSize: 18, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>-</button>
                  <span style={{ minWidth: 36, textAlign: 'center', color: '#fff', fontSize: 14 }}>{Math.round(previewScale * 100)}%</span>
                  <button onClick={handleZoomIn} style={{ width: 28, height: 28, fontSize: 18, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>+</button>
                  <button onClick={handleResetZoom} style={{ width: 28, height: 28, fontSize: 12, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>1:1</button>
                  <button onClick={handleFitToScreen} style={{ width: 28, height: 28, fontSize: 12, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>适应</button>
                  {error && (
                    <div style={{ 
                      marginLeft: 8, 
                      padding: '4px 8px', 
                      background: errorType === 'network' ? '#ff6b6b' : '#f44747', 
                      color: '#fff', 
                      borderRadius: 4, 
                      fontSize: 12, 
                      display: 'flex', 
                      alignItems: 'center', 
                      gap: 4,
                      maxWidth: 200,
                      whiteSpace: 'nowrap',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis'
                    }}>
                      <span style={{ 
                        overflow: 'hidden', 
                        textOverflow: 'ellipsis',
                        fontSize: 11
                      }}>{error}</span>
                      <button
                        onClick={() => { setError(''); setErrorType(''); }}
                        style={{
                          background: 'none',
                          border: 'none',
                          color: '#fff',
                          cursor: 'pointer',
                          padding: 0,
                          fontSize: 14,
                          lineHeight: 1,
                          opacity: 0.8
                        }}
                        onMouseEnter={(e) => e.currentTarget.style.opacity = '1'}
                        onMouseLeave={(e) => e.currentTarget.style.opacity = '0.8'}
                      >
                        ×
                      </button>
                    </div>
                  )}
                </div>
              </div>
              <div
                ref={previewRef}
                className="ope-svg-preview"
                style={{ flex: 1, minHeight: 0, background: '#fff', overflow: 'hidden', position: 'relative' }}
              >
                <div 
                  style={{ 
                    transform: `scale(${previewScale}) translate(${panPosition.x}px, ${panPosition.y}px)`, 
                    transformOrigin: 'center center', 
                    transition: isDragging ? 'none' : 'transform 0.2s', 
                    display: 'inline-block',
                    minWidth: '100%',
                    minHeight: '100%',
                    cursor: isDragging ? 'grabbing' : 'grab',
                    userSelect: 'none'
                  }}
                  onMouseDown={handleMouseDown}
                  dangerouslySetInnerHTML={{ __html: svgContent }}
                />
              </div>
            </div>
          )}
          {mode === 'split' && (
            <div className="ope-split-row" style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              width: '100%',
              height: '100%',
              minWidth: 0,
              minHeight: 0,
              boxSizing: 'border-box',
              flex: 1,
              overflow: 'hidden',
            }}>
              {/* 左侧代码区 */}
              <div style={{
                minWidth: 0,
                minHeight: 0,
                height: '100%',
                width: '100%',
                display: 'flex',
                flexDirection: 'column',
                boxSizing: 'border-box',
                overflow: 'hidden',
              }}>
                <div style={{
                  height: 48,
                  display: 'flex',
                  alignItems: 'center',
                  padding: '0 16px',
                  fontWeight: 600,
                  fontSize: 16,
                  background: '#222',
                  color: '#fff',
                  borderTopLeftRadius: 8,
                  boxSizing: 'border-box',
                  flexShrink: 0,
                }}>
                  <span style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>PUML 代码</span>
                </div>
                <div className="ope-textarea-container" style={{ flex: 1, position: 'relative' }}>
                  {/* 错误行高亮 */}
                  {errorLinePositions.map((pos: any, index: number) => (
                    <React.Fragment key={`error-${pos.lineNumber}-${index}`}>
                      <div
                        className="ope-error-highlight"
                        style={{
                          top: pos.top,
                          height: pos.height,
                        }}
                      />
                      <div
                        className="ope-error-underline"
                        style={{
                          top: pos.top + pos.height - 2,
                        }}
                      />
                    </React.Fragment>
                  ))}
                  <textarea
                    ref={editorRef}
                    value={isControlled ? value : pumlCode}
                    onChange={handleInputChange}
                    className="ope-textarea"
                    placeholder="输入 PlantUML 代码..."
                    readOnly={readOnly}
                    style={{
                      flex: 1,
                      width: '100%',
                      height: '100%',
                      minHeight: 0,
                      background: '#222',
                      color: '#fff',
                      border: 'none',
                      borderBottomLeftRadius: 8,
                      fontSize: 15,
                      fontFamily: 'Fira Mono, Menlo, Monaco, Consolas, monospace',
                      boxSizing: 'border-box',
                      resize: 'none',
                      padding: 16,
                      overflow: 'auto',
                      position: 'relative',
                      zIndex: 3,
                    }}
                  />
                </div>
              </div>
              {/* 右侧预览区 */}
              <div style={{
                minWidth: 0,
                minHeight: 0,
                height: '100%',
                width: '100%',
                display: 'flex',
                flexDirection: 'column',
                boxSizing: 'border-box',
                overflow: 'hidden',
              }}>
                <div style={{
                  height: 48,
                  display: 'flex',
                  alignItems: 'center',
                  padding: '0 16px',
                  fontWeight: 600,
                  fontSize: 16,
                  background: '#222',
                  color: '#fff',
                  borderTopRightRadius: 8,
                  boxSizing: 'border-box',
                  flexShrink: 0,
                  position: 'relative',
                }}>
                  <span style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>预览</span>
                  <div style={{ position: 'absolute', right: 16, top: '50%', transform: 'translateY(-50%)', display: 'flex', alignItems: 'center', gap: 4 }}>
                    <button onClick={handleZoomOut} style={{ width: 28, height: 28, fontSize: 18, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>-</button>
                    <span style={{ minWidth: 36, textAlign: 'center', color: '#fff', fontSize: 14 }}>{Math.round(previewScale * 100)}%</span>
                    <button onClick={handleZoomIn} style={{ width: 28, height: 28, fontSize: 18, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>+</button>
                    <button onClick={handleResetZoom} style={{ width: 28, height: 28, fontSize: 12, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>1:1</button>
                    <button onClick={handleFitToScreen} style={{ width: 28, height: 28, fontSize: 12, borderRadius: 4, border: 'none', background: '#444', color: '#fff', cursor: 'pointer' }}>适应</button>
                    {error && (
                      <div style={{ 
                        marginLeft: 8, 
                        padding: '4px 8px', 
                        background: errorType === 'network' ? '#ff6b6b' : '#f44747', 
                        color: '#fff', 
                        borderRadius: 4, 
                        fontSize: 12, 
                        display: 'flex', 
                        alignItems: 'center', 
                        gap: 4,
                        maxWidth: 200,
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis'
                      }}>
                        <span style={{ 
                          overflow: 'hidden', 
                          textOverflow: 'ellipsis',
                          fontSize: 11
                        }}>{error}</span>
                        <button
                          onClick={() => { setError(''); setErrorType(''); }}
                          style={{
                            background: 'none',
                            border: 'none',
                            color: '#fff',
                            cursor: 'pointer',
                            padding: 0,
                            fontSize: 14,
                            lineHeight: 1,
                            opacity: 0.8
                          }}
                          onMouseEnter={(e) => e.currentTarget.style.opacity = '1'}
                          onMouseLeave={(e) => e.currentTarget.style.opacity = '0.8'}
                        >
                          ×
                        </button>
                      </div>
                    )}
                  </div>
                </div>
                <div
                  ref={previewRef}
                  className="ope-svg-preview"
                  style={{
                    flex: 1,
                    width: '100%',
                    height: '100%',
                    minHeight: 0,
                    overflow: 'hidden',
                    background: '#fff',
                    borderBottomRightRadius: 8,
                    boxSizing: 'border-box',
                    position: 'relative'
                  }}
                >
                  <div 
                    style={{ 
                      transform: `scale(${previewScale}) translate(${panPosition.x}px, ${panPosition.y}px)`, 
                      transformOrigin: 'center center', 
                      transition: isDragging ? 'none' : 'transform 0.2s', 
                      display: 'inline-block',
                      minWidth: '100%',
                      minHeight: '100%',
                      cursor: isDragging ? 'grabbing' : 'grab',
                      userSelect: 'none'
                    }}
                    onMouseDown={handleMouseDown}
                    dangerouslySetInnerHTML={{ __html: svgContent }}
                  />
                </div>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  );
};

export default OnlinePUMLEditor;