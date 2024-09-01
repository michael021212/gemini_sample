package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	return f1(ctx)
}

func printCandidates(cs []*genai.Candidate) {
	for _, c := range cs {
		for _, p := range c.Content.Parts {
			fmt.Println((p))
		}
	}
}

// テキストのみの入力からテキストを生成する
// https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#generate-text-from-text
func f1(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	prompt := genai.Text("ONEPIECEの作者は？")

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return err
	}

	printCandidates(resp.Candidates)

	return nil
}

// テキストと画像の入力からテキストを生成する
// https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#generate-text-from-text-and-image
func f2(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	imgData1, err := os.ReadFile("onepiece01_luffy.png")
	if err != nil {
		return err
	}

	imgData2, err := os.ReadFile("onepiece02_zoro_bandana.png")
	if err != nil {
		return err
	}

	prompt := []genai.Part{
		genai.ImageData("png", imgData1),
		genai.ImageData("png", imgData2),
		genai.Text("２つの画像の違いを教えて"),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return err
	}

	printCandidates(resp.Candidates)

	return nil
}

// インタラクティブなチャットを作成する
// https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#chat
func f3(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	cs := model.StartChat()

	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("こんにちは、赤い果物は何がありますか？"),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("いちごやりんごがあります。"),
			},
			Role: "model",
		},
	}

	resp, err := cs.SendMessage(ctx, genai.Text("その中で小さいのはどちらですか？"))
	if err != nil {
		return err
	}

	printCandidates(resp.Candidates)

	return nil
}

// テキスト ストリームを生成する(結果全体を待たない)
// https://ai.google.dev/gemini-api/docs/text-generation?hl=ja&lang=go#generate-a-text-stream
func f4(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	prompt := genai.Text("モンキー・D・ルフィが食べた悪魔の実について説明して")

	iter := model.GenerateContentStream(ctx, prompt)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		printCandidates(resp.Candidates)
	}

	return nil
}
