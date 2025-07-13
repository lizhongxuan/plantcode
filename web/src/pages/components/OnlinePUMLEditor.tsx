import React, { useState, useEffect, useRef, useCallback } from 'react';
import './OnlinePUMLEditor.css';

interface OnlinePUMLEditorProps {
  initialCode: string;
  onClose: () => void;
  onSave: (newCode: string) => void;
}

const OnlinePUMLEditor: React.FC<OnlinePUMLEditorProps> = ({ initialCode, onClose, onSave }) => {
  const [pumlCode, setPumlCode] = useState(initialCode);
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
    onSave(pumlCode);
    onClose();
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

        const lines = editor.value.split('\\n');
        if (lineNumber > 0 && lineNumber <= lines.length) {
            let charPosition = 0;
            for (let i = 0; i < lineNumber - 1; i++) {
                charPosition += lines[i].length + 1; // +1 for newline
            }
            
            editor.focus();
            editor.setSelectionRange(charPosition, charPosition + lines[lineNumber - 1].length);

            // Scroll to the selection
            const lineHeight = parseFloat(getComputedStyle(editor).lineHeight);
            const scrollTop = Math.max(0, (lineNumber - 5) * lineHeight);
            editor.scrollTop = scrollTop;
        }
    };

    const jumpToCodeElement = (elementName: string) => {
        const editor = editorRef.current;
        if (!editor) return;

        const code = editor.value;
        const lines = code.split('\\n');
        const searchName = elementName.toLowerCase().trim();

        if (!searchName) return;

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i].toLowerCase();
            if (line.includes(searchName)) {
                jumpToLine(i + 1);
                return;
            }
        }
    };

    const extractElementName = (element: Element): string => {
        let name = element.textContent?.trim() || '';
        if (!name && element.querySelector('text')) {
            name = element.querySelector('text')!.textContent?.trim() || '';
        }
        return name;
    };

    const previewEl = previewRef.current;
    if (!previewEl || !svgContent) return;
    
    // Clear previous content and listeners by re-rendering
    previewEl.innerHTML = svgContent;

    const svg = previewEl.querySelector('svg');
    if (!svg) return;

    const clickableElements = svg.querySelectorAll('text, rect, ellipse, circle, polygon, path, g');
    
    const clickHandler = (event: Event) => {
        event.preventDefault();
        event.stopPropagation();
        
        const elementName = extractElementName(event.currentTarget as Element);
        if (elementName) {
            jumpToCodeElement(elementName);
        }
    };

    const mouseEnterHandler = (event: Event) => {
      const target = event.currentTarget as SVGElement;
      target.style.opacity = '0.7';
    };

    const mouseLeaveHandler = (event: Event) => {
      const target = event.currentTarget as SVGElement;
      target.style.opacity = '1';
    };

    clickableElements.forEach(element => {
        (element as HTMLElement).style.cursor = 'pointer';
        element.addEventListener('click', clickHandler);
        element.addEventListener('mouseenter', mouseEnterHandler);
        element.addEventListener('mouseleave', mouseLeaveHandler);
    });

    return () => {
      clickableElements.forEach(element => {
        element.removeEventListener('click', clickHandler);
        element.removeEventListener('mouseenter', mouseEnterHandler);
        element.removeEventListener('mouseleave', mouseLeaveHandler);
      });
    };
  }, [svgContent]);

  return (
    <div className="online-puml-editor-modal">
      <div className="online-puml-editor-container">
        <header className="ope-header">
          <h1>PlantUML 在线编辑器</h1>
          <div>
            <button onClick={handleSave} className="ope-button ope-save-btn">保存并关闭</button>
            <button onClick={onClose} className="ope-close-btn">&times;</button>
          </div>
        </header>
        <main className="ope-main">
          <div className="ope-panel ope-editor-panel">
            <div className="ope-panel-header">
              <h2>PUML 代码</h2>
              <button onClick={handleRender} disabled={isLoading} className="ope-button">
                {isLoading ? '渲染中...' : '渲染图表 (Ctrl+Enter)'}
              </button>
            </div>
            <textarea
              ref={editorRef}
              value={pumlCode}
              onChange={(e) => setPumlCode(e.target.value)}
              className="ope-textarea"
              placeholder="输入 PlantUML 代码..."
            />
          </div>
          <div className="ope-panel ope-preview-panel">
            <div className="ope-panel-header">
              <h2>预览</h2>
            </div>
            {error && <div className="ope-error-display">{error}</div>}
            <div
              ref={previewRef}
              className="ope-svg-preview"
              dangerouslySetInnerHTML={{ __html: svgContent }}
            />
          </div>
        </main>
      </div>
    </div>
  );
};

export default OnlinePUMLEditor;