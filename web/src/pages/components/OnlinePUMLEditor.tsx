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
  const editorRef = useRef<HTMLTextAreaElement>(null);
  const previewRef = useRef<HTMLDivElement>(null);

  const handleRender = useCallback(async () => {
    setIsLoading(true);
    setError('');
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
        setError(errorMessage);
        setSvgContent('');
      }
    } catch (err) {
      setError('请求渲染服务失败');
      setSvgContent('');
    } finally {
      setIsLoading(false);
    }
  }, [pumlCode]);

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
          onClick={error ? (e) => {
            // 只在点击非错误提示框区域时关闭
            if ((e.target as HTMLElement).closest('.ope-error-display') == null) setError('');
          } : undefined}
        >
          {error && (
            <div className="ope-error-display error-with-close" style={{ position: 'absolute', top: 32, left: '50%', transform: 'translateX(-50%)', zIndex: 20, maxWidth: '80%', minWidth: 320 }} onClick={e => e.stopPropagation()}>
              <button
                className="ope-error-close-btn"
                onClick={() => setError('')}
                aria-label="关闭错误提示"
              >
                ×
              </button>
              {error}
            </div>
          )}
          {mode === 'edit' && (
            <div className="ope-panel ope-editor-panel" style={{ flex: 1, minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <div className="ope-panel-header">
                <h2>PUML 代码</h2>
              </div>
              <textarea
                ref={editorRef}
                value={isControlled ? value : pumlCode}
                onChange={handleInputChange}
                className="ope-textarea"
                placeholder="输入 PlantUML 代码..."
                readOnly={readOnly}
                style={{ flex: 1, minHeight: 0, resize: 'none', background: '#fff' }}
              />
            </div>
          )}
          {mode === 'preview' && (
            <div className="ope-panel ope-preview-panel" style={{ flex: 1, minHeight: 0, background: 'transparent', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <div className="ope-panel-header">
                <h2>预览</h2>
              </div>
              <div
                ref={previewRef}
                className="ope-svg-preview"
                style={{ flex: 1, minHeight: 0, background: '#fff', overflow: 'auto' }}
                dangerouslySetInnerHTML={{ __html: svgContent }}
              />
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
                  }}
                />
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
                }}>
                  <span style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>预览</span>
                </div>
                <div
                  ref={previewRef}
                  className="ope-svg-preview"
                  style={{
                    flex: 1,
                    width: '100%',
                    height: '100%',
                    minHeight: 0,
                    overflow: 'auto',
                    background: '#fff',
                    borderBottomRightRadius: 8,
                    boxSizing: 'border-box',
                    padding: 16,
                  }}
                  dangerouslySetInnerHTML={{ __html: svgContent }}
                />
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  );
};

export default OnlinePUMLEditor;