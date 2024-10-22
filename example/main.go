package main

import (
	"fmt"
	"math"

	"github.com/kelindar/llm"
)

func main() {
	//exampleCompleteText()
	exampleEmbedText()
}

func exampleCompleteText() {
	m, err := llm.New("../dist/Llama-3.2-1B-Instruct-Q6_K_L.gguf", 0)
	//m, err := llm.New("../dist/MiniLM-L6-v2.Q4_K_M.gguf", 0)
	if err != nil {
		panic(err)
	}

	defer m.Close()

	ctx := m.Context(0)
	defer ctx.Close()

	template := `### System:
	You are a character in an adventure game.

	### Instruction:
	%s

	### Response (10 words, engaging, natural, authentic, descriptive, creative):
	`

	for {
		var input string
		fmt.Printf("\n >> ")
		fmt.Scanln(&input)

		text := fmt.Sprintf(template, input)
		out, err := ctx.CompleteText(text, 2048)
		if err != nil {
			panic(err)
		}

		fmt.Println(out)
	}

}

func exampleEmbedText() {
	m, err := llm.New("../dist/MiniLM-L6-v2.Q4_K_M.gguf", 0)
	//m, err = llm.New("../dist/e5-small-v2.Q4_K_M.gguf", 0)
	//m, err := llm.New("../dist/e5-base-v2.Q5_K_M.gguf", 0)
	if err != nil {
		panic(err)
	}

	defer m.Close()

	prompts := []string{
		"The quick brown fox jumps over the lazy dog.",
		"A fast, nimble fox leaps over a sleepy dog.",
		"An agile brown llama hops over a tired pug.",
		"A beautiful sunset on the beach paints the sky with vibrant colors.",
	}

	embeddings := make([][]float32, len(prompts))
	for i, prompt := range prompts {
		embeddings[i], err = m.EmbedText(prompt)
		if err != nil {
			panic(err)
		}
	}

	// Compute pairwise cosine similarities and print them out
	for i := 0; i < len(embeddings); i++ {
		for j := i + 1; j < len(embeddings); j++ {
			cos := cosine(embeddings[i], embeddings[j])
			fmt.Printf("\n * Similarity = %.2f\n", cos)
			fmt.Printf("   1: %s\n", prompts[i])
			fmt.Printf("   2: %s\n", prompts[j])
		}
	}
}

// cosine computes the cosine similarity between two vectors. Higher values
// indicate more similar vectors.
func cosine(a, b []float32) float64 {
	if len(a) != len(b) {
		panic("vectors must be of equal length")
	}

	dp, an, bn := float64(0), float64(0), float64(0)
	for i := range a {
		dp += float64(a[i] * b[i])
		an += float64(a[i] * a[i])
		bn += float64(b[i] * b[i])
	}

	return dp / (math.Sqrt(an) * math.Sqrt(bn))
}
