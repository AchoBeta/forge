/**
 * Forge Mind Map Frontend Application
 * 基于生成式AI的思维导图系统前端
 */

// API Configuration
const API_BASE_URL = 'http://localhost:8000';

// DOM Elements
const topicInput = document.getElementById('topic');
const depthInput = document.getElementById('depth');
const branchesInput = document.getElementById('branches');
const generateBtn = document.getElementById('generateBtn');
const mindmapDiv = document.getElementById('mindmap');

// State
let currentMindMap = null;

/**
 * Initialize the application
 */
function init() {
    generateBtn.addEventListener('click', generateMindMap);
    
    // Allow Enter key to trigger generation
    topicInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            generateMindMap();
        }
    });
    
    // Generate initial mind map
    generateMindMap();
}

/**
 * Generate mind map from API
 */
async function generateMindMap() {
    const topic = topicInput.value.trim();
    const depth = parseInt(depthInput.value);
    const branches = parseInt(branchesInput.value);
    
    if (!topic) {
        alert('请输入主题！');
        return;
    }
    
    // Update button state
    setLoadingState(true);
    
    try {
        const response = await fetch(`${API_BASE_URL}/generate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                topic: topic,
                depth: depth,
                branches: branches
            })
        });
        
        if (!response.ok) {
            throw new Error(`API Error: ${response.status}`);
        }
        
        const data = await response.json();
        currentMindMap = data;
        
        // Render mind map
        renderMindMap(data.root);
        
    } catch (error) {
        console.error('Error generating mind map:', error);
        alert('生成思维导图时出错，请检查后端服务是否运行在 http://localhost:8000');
    } finally {
        setLoadingState(false);
    }
}

/**
 * Set loading state for generate button
 */
function setLoadingState(loading) {
    const btnText = generateBtn.querySelector('.btn-text');
    const btnLoading = generateBtn.querySelector('.btn-loading');
    
    if (loading) {
        btnText.style.display = 'none';
        btnLoading.style.display = 'inline';
        generateBtn.disabled = true;
    } else {
        btnText.style.display = 'inline';
        btnLoading.style.display = 'none';
        generateBtn.disabled = false;
    }
}

/**
 * Render mind map visualization
 */
function renderMindMap(rootNode) {
    // Clear existing content
    mindmapDiv.innerHTML = '';
    
    // Calculate positions for all nodes
    const positions = calculateNodePositions(rootNode);
    
    // Render connections first (so they appear behind nodes)
    renderConnections(positions);
    
    // Render nodes
    renderNodes(positions);
}

/**
 * Calculate positions for all nodes in the mind map
 */
function calculateNodePositions(rootNode) {
    const positions = [];
    const container = mindmapDiv;
    const width = container.clientWidth;
    const height = container.clientHeight;
    
    // Root node in center
    const rootX = width / 2;
    const rootY = height / 2;
    
    function traverse(node, x, y, angle, radius, level, parent) {
        const pos = {
            node: node,
            x: x,
            y: y,
            level: level,
            parent: parent
        };
        positions.push(pos);
        
        if (node.children && node.children.length > 0) {
            const childCount = node.children.length;
            const angleStep = (Math.PI * 2) / childCount;
            const nextRadius = radius * 0.7;
            
            node.children.forEach((child, index) => {
                const childAngle = angle + (index - (childCount - 1) / 2) * angleStep;
                const childX = x + Math.cos(childAngle) * radius;
                const childY = y + Math.sin(childAngle) * radius;
                traverse(child, childX, childY, childAngle, nextRadius, level + 1, pos);
            });
        }
    }
    
    traverse(rootNode, rootX, rootY, 0, 200, 0, null);
    return positions;
}

/**
 * Render connection lines between nodes
 */
function renderConnections(positions) {
    positions.forEach(pos => {
        if (pos.parent) {
            const connection = document.createElement('div');
            connection.className = 'connection';
            
            const dx = pos.x - pos.parent.x;
            const dy = pos.y - pos.parent.y;
            const distance = Math.sqrt(dx * dx + dy * dy);
            const angle = Math.atan2(dy, dx);
            
            connection.style.width = `${distance}px`;
            connection.style.left = `${pos.parent.x}px`;
            connection.style.top = `${pos.parent.y}px`;
            connection.style.transform = `rotate(${angle}rad)`;
            
            mindmapDiv.appendChild(connection);
        }
    });
}

/**
 * Render node elements
 */
function renderNodes(positions) {
    positions.forEach(pos => {
        const nodeElement = document.createElement('div');
        nodeElement.className = `node level-${pos.level}`;
        nodeElement.textContent = pos.node.text;
        nodeElement.title = pos.node.text; // Tooltip
        
        // Position the node (centered on the coordinate)
        nodeElement.style.left = `${pos.x}px`;
        nodeElement.style.top = `${pos.y}px`;
        nodeElement.style.transform = 'translate(-50%, -50%)';
        
        // Add click handler
        nodeElement.addEventListener('click', () => {
            console.log('Node clicked:', pos.node);
            // Could expand/collapse or show details
        });
        
        mindmapDiv.appendChild(nodeElement);
    });
}

// Initialize application when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}
