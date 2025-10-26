"""
FastAPI backend for AI-powered mind map generation system
基于生成式AI的思维导图系统后端
"""
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import json
import os

app = FastAPI(
    title="Forge Mind Map API",
    description="AI-powered mind map generation system",
    version="1.0.0"
)

# Enable CORS for frontend communication
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class MindMapNode(BaseModel):
    """Mind map node structure"""
    id: str
    text: str
    children: Optional[List['MindMapNode']] = []
    level: int = 0


class GenerateRequest(BaseModel):
    """Request model for mind map generation"""
    topic: str
    depth: int = 3
    branches: int = 3


class MindMapResponse(BaseModel):
    """Response model for mind map"""
    root: MindMapNode
    metadata: dict


def generate_mindmap_with_ai(topic: str, depth: int = 3, branches: int = 3) -> MindMapNode:
    """
    Generate mind map using AI logic
    This is a simplified implementation that can be extended with actual AI models
    """
    
    # AI-powered mind map generation logic
    # In a production system, this would call OpenAI API, local LLM, etc.
    
    def create_node(text: str, level: int, node_id: str) -> MindMapNode:
        node = MindMapNode(id=node_id, text=text, level=level, children=[])
        
        if level < depth:
            # Generate child nodes based on topic context
            children_topics = generate_subtopics(text, level, branches)
            for i, child_topic in enumerate(children_topics):
                child_id = f"{node_id}_{i}"
                child_node = create_node(child_topic, level + 1, child_id)
                node.children.append(child_node)
        
        return node
    
    # Create root node
    root = create_node(topic, 0, "root")
    return root


def generate_subtopics(parent_topic: str, level: int, count: int) -> List[str]:
    """
    Generate subtopics based on the parent topic
    This simulates AI-generated content
    """
    # AI-powered subtopic generation
    # This is a simplified version - in production, use actual AI models
    
    topic_templates = {
        0: {  # First level
            "default": [
                f"{parent_topic}的核心概念",
                f"{parent_topic}的应用场景",
                f"{parent_topic}的发展趋势",
                f"{parent_topic}的关键技术",
                f"{parent_topic}的优势与挑战"
            ]
        },
        1: {  # Second level
            "default": [
                "主要特点",
                "实现方法",
                "最佳实践",
                "常见问题",
                "解决方案"
            ]
        },
        2: {  # Third level
            "default": [
                "详细说明",
                "技术细节",
                "案例分析",
                "注意事项"
            ]
        }
    }
    
    templates = topic_templates.get(level, topic_templates[2])["default"]
    return templates[:count]


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Forge Mind Map API - 基于生成式AI的思维导图系统",
        "version": "1.0.0",
        "endpoints": {
            "/generate": "POST - Generate mind map from topic",
            "/health": "GET - Health check"
        }
    }


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "forge-mindmap"}


@app.post("/generate", response_model=MindMapResponse)
async def generate_mindmap(request: GenerateRequest):
    """
    Generate mind map from topic using AI
    """
    try:
        # Validate input
        if not request.topic or len(request.topic.strip()) == 0:
            raise HTTPException(status_code=400, detail="Topic cannot be empty")
        
        if request.depth < 1 or request.depth > 5:
            raise HTTPException(status_code=400, detail="Depth must be between 1 and 5")
        
        if request.branches < 1 or request.branches > 8:
            raise HTTPException(status_code=400, detail="Branches must be between 1 and 8")
        
        # Generate mind map
        root_node = generate_mindmap_with_ai(
            request.topic,
            request.depth,
            request.branches
        )
        
        # Create response
        response = MindMapResponse(
            root=root_node,
            metadata={
                "topic": request.topic,
                "depth": request.depth,
                "branches": request.branches,
                "generated_by": "Forge AI System"
            }
        )
        
        return response
        
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error generating mind map: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    port = int(os.environ.get("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
