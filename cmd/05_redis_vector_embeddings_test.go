package main

/*

*/

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/knights-analytics/hugot"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisVectorEmbeddingsTest(t *testing.T) {
	t.Skip("Lock problems during test")

	rdb := setupTestRedisClient()
	defer rdb.Close()
	ctx := context.Background()

	err := flushRedisDB(ctx, rdb)
	require.NoError(t, err)

	session, err := hugot.NewGoSession()
	require.NoError(t, err)

	defer func() {
		err := session.Destroy()
		require.NoError(t, err)
	}()

	downloadOptions := hugot.NewDownloadOptions()
	downloadOptions.OnnxFilePath = "onnx/model.onnx"
	modelPath, err := hugot.DownloadModel(
		"sentence-transformers/all-MiniLM-L6-v2",
		"./models/",
		downloadOptions,
	)
	require.NoError(t, err)

	config := hugot.FeatureExtractionConfig{
		ModelPath: modelPath,
		Name:      "embeddingPipeline",
	}

	embeddingPipeline, err := hugot.NewPipeline(session, config)
	require.NoError(t, err)

	for name, details := range peopleData {
		result, err := embeddingPipeline.RunPipeline([]string{details.Description})
		require.NoError(t, err)

		embFloat32 := result.Embeddings[0]
		embFloat64 := make([]float64, len(embFloat32))
		for i, v := range embFloat32 {
			embFloat64[i] = float64(v)
		}

		_, err = rdb.VAdd(ctx, "famousPeople", name, &redis.VectorValues{Val: embFloat64}).Result()
		require.NoError(t, err)

		_, err = rdb.VSetAttr(ctx, "famousPeople", name, map[string]interface{}{
			"born": details.Born,
			"died": details.Died,
		}).Result()
		require.NoError(t, err)
	}

	queryValue := "actors"
	queryResult, err := embeddingPipeline.RunPipeline([]string{queryValue})
	require.NoError(t, err)

	queryFloat32 := queryResult.Embeddings[0]
	queryFloat64 := make([]float64, len(queryFloat32))
	for i, v := range queryFloat32 {
		queryFloat64[i] = float64(v)
	}

	actorsResults, err := rdb.VSim(ctx, "famousPeople", &redis.VectorValues{Val: queryFloat64}).Result()
	require.NoError(t, err)

	fmt.Printf("'actors': %v\n", strings.Join(actorsResults, ", "))
}

type PersonData struct {
	Born        int
	Died        int
	Description string
}

var peopleData = map[string]PersonData{
	"Marie Curie": {
		Born: 1867, Died: 1934,
		Description: `Polish-French chemist and physicist. The only person ever to win
		two Nobel prizes for two different sciences.
		`,
	},
	"Linus Pauling": {
		Born: 1901, Died: 1994,
		Description: `American chemist and peace activist. One of only two people to win two
		Nobel prizes in different fields (chemistry and peace).
		`,
	},
	"Freddie Mercury": {
		Born: 1946, Died: 1991,
		Description: `British musician, best known as the lead singer of the rock band
		Queen.
		`,
	},
	"Marie Fredriksson": {
		Born: 1958, Died: 2019,
		Description: `Swedish multi-instrumentalist, mainly known as the lead singer and
		keyboardist of the band Roxette.
		`,
	},
	"Paul Erdos": {
		Born: 1913, Died: 1996,
		Description: `Hungarian mathematician, known for his eccentric personality almost
		as much as his contributions to many different fields of mathematics.
		`,
	},
	"Maryam Mirzakhani": {
		Born: 1977, Died: 2017,
		Description: `Iranian mathematician. The first woman ever to win the Fields medal
		for her contributions to mathematics.
		`,
	},
	"Masako Natsume": {
		Born: 1957, Died: 1985,
		Description: `Japanese actress. She was very famous in Japan but was primarily
		known elsewhere in the world for her portrayal of Tripitaka in the
		TV series Monkey.
		`,
	},
	"Chaim Topol": {
		Born: 1935, Died: 2023,
		Description: `Israeli actor and singer, usually credited simply as 'Topol'. He was
		best known for his many appearances as Tevye in the musical Fiddler
		on the Roof.
		`,
	},
}
