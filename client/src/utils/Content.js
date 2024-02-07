function mapCodeBlocksToPre(text) {
    const codeBlockRegex = /```([\s\S]+?)```/g; // Regular expression to match code blocks
    
    return text.replace(codeBlockRegex, (match, code) => {
        return `<pre><code>${code}</code></pre>`;
    });
}

export{
    mapCodeBlocksToPre
}