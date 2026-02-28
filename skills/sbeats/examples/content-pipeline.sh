#!/bin/bash
# Example: Content Creation Workflow
# This example shows how to manage a content pipeline

set -e

echo "=== Content Creation Workflow Example ==="
echo ""

sb init

echo "Creating content pipeline..."
echo ""

# Main content strategy
echo "=== Content Strategy ==="
sb create "Q1 Content Calendar" -d "Plan blog posts, videos, and social content for Q1" -p 0
# sb-strategy

echo ""
echo "=== Blog Post: Getting Started with Go ==="
sb create "Research Go best practices" -d "Find latest patterns and tools" -p 1 --parent sb-strategy
# sb-research1

sb create "Write blog outline" -d "Structure: intro, setup, hello world, packages, conclusion" -p 1 --parent sb-strategy --deps sb-research1
# sb-outline1

sb create "Write draft content" -d "Write 1500 words with code examples" -p 2 --parent sb-strategy --deps sb-outline1
# sb-draft1

sb create "Create code examples" -d "Build working example projects" -p 2 --parent sb-strategy --deps sb-outline1
# sb-code1

sb create "Review and edit" -d "Technical review and copy editing" -p 2 --parent sb-strategy --deps sb-draft1
# sb-review1

sb create "Create header image" -d "Design blog post header in Canva" -p 3 --parent sb-strategy --deps sb-outline1
# sb-image1

sb create "Publish blog post" -d "Upload to CMS, schedule publish" -p 1 --parent sb-strategy --deps sb-review1,sb-code1,sb-image1
# sb-publish1

echo ""
echo "=== Video Tutorial: Docker Basics ==="
sb create "Write video script" -d "Plan video sections and talking points" -p 1 --parent sb-strategy
# sb-script2

sb create "Record video" -d "Screen recording with voiceover" -p 2 --parent sb-strategy --deps sb-script2
# sb-record2

sb create "Edit video" -d "Add intros, outros, captions" -p 2 --parent sb-strategy --deps sb-record2
# sb-edit2

sb create "Create thumbnail" -d "Design YouTube thumbnail" -p 3 --parent sb-strategy --deps sb-script2
# sb-thumb2

sb create "Publish video" -d "Upload to YouTube with SEO metadata" -p 1 --parent sb-strategy --deps sb-edit2,sb-thumb2
# sb-publish2

echo ""
echo "=== Social Media ==="
sb create "Create Twitter thread" -d "Summarize blog post in 5 tweets" -p 2 --parent sb-strategy --deps sb-publish1

sb create "Create LinkedIn post" -d "Professional version for LinkedIn" -p 2 --parent sb-strategy --deps sb-publish1

sb create "Schedule social posts" -d "Use Buffer to schedule posts" -p 2 --parent sb-strategy

echo ""
echo "=== Content Pipeline Status ==="
echo ""
echo "Ready to work on:"
sb ready | head -8

echo ""
echo "View content calendar:"
sb list --parent sb-strategy --all

echo ""
echo "Content workflow tips:"
echo "- Research → Outline → Draft → Review → Publish"
echo "- Use dependencies to ensure proper order"
echo "- Parallel work possible (writing + images)"
echo "- Review sb ready daily for next content piece"
