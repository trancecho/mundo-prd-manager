package api

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"log"
)

// GenerateBlock 创建飞书文档并写入内容
func GenerateBlock(ctx context.Context, folderToken string, title string, content string) (string, error) {
	// 创建文档
	response, err := Cli.Client.Docx.Document.Create(ctx, larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken(folderToken). // 文件夹 token，传空表示在根目录创建文档
			Title(title).             // 文档标题
			Build()).
		Build())

	log.Println(response)
	if err != nil {
		return "", err
	}

	if response == nil || response.Data == nil {
		return "", fmt.Errorf("创建文档API返回结果为空")
	}
	docToken := *response.Data.Document.DocumentId
	revisionId := response.Data.Document.RevisionId
	if revisionId == nil {
		// 如果没有返回修订ID，使用默认值0
		defaultRevisionId := 0
		revisionId = &defaultRevisionId
	}

	req := larkdocx.NewCreateDocumentBlockChildrenReqBuilder().
		DocumentId(docToken).
		BlockId(docToken).
		DocumentRevisionId(*revisionId).
		Body(larkdocx.NewCreateDocumentBlockChildrenReqBodyBuilder().
			Children([]*larkdocx.Block{
				larkdocx.NewBlockBuilder().
					BlockType(2).
					Text(larkdocx.NewTextBuilder().
						Style(larkdocx.NewTextStyleBuilder().
							Build()).
						Elements([]*larkdocx.TextElement{
							larkdocx.NewTextElementBuilder().
								TextRun(larkdocx.NewTextRunBuilder().
									Content(content).
									TextElementStyle(larkdocx.NewTextElementStyleBuilder().
										Bold(true).
										BackgroundColor(14).
										TextColor(5).
										Build()).
									Build()).
								Build(),
						}).
						Build()).
					Build(),
			}).
			Index(0).
			Build()).
		Build()

	log.Println("创建文档块中···")
	resp, err := Cli.Client.Docx.V1.DocumentBlockChildren.Create(ctx, req)
	if err != nil {
		return "", err
	}

	//服务端处理错误
	if !resp.Success() {
		log.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", fmt.Errorf("failed to create document block children: %s", resp.CodeError)
	}

	return docToken, nil
}
