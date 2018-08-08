package sqlsol_test

import (
	"testing"

	"github.com/hiturria/bosmarmot/vent/sqlsol"
	"github.com/hiturria/bosmarmot/vent/types"
	"github.com/stretchr/testify/require"
)

func TestSetBlockID(t *testing.T) {
	t.Run("successfully sets an id block", func(t *testing.T) {
		blockData := sqlsol.NewBlockData()
		blockData.SetBlockID("44")

		blockID := blockData.GetBlockID()
		require.Equal(t, "44", blockID)
	})
}

func TestGetBlockID(t *testing.T) {
	t.Run("successfully gets an id block", func(t *testing.T) {
		blockData := sqlsol.NewBlockData()
		blockData.SetBlockID("99")

		blockID := blockData.GetBlockID()
		require.Equal(t, "99", blockID)
	})
}

func TestSetRow(t *testing.T) {
	t.Run("successfully sets a new data row", func(t *testing.T) {
		row := make(types.EventDataRow)
		row["c1"] = "v1"
		row["c2"] = "v2"

		blockData := sqlsol.NewBlockData()
		blockData.AddRow("TEST_TABLE", row)

		rows, err := blockData.GetRows("TEST_TABLE")
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
		require.Equal(t, "v1", rows[0]["c1"])
		require.Equal(t, "v2", rows[0]["c2"])
	})
}

func TestGetBlockData(t *testing.T) {
	t.Run("successfully gets block data", func(t *testing.T) {
		blockData := sqlsol.NewBlockData()
		blk := blockData.GetBlockData()
		require.Equal(t, "", blk.Block)
	})
}

func TestPendingRows(t *testing.T) {
	t.Run("successfully returns true if a given block has pending rows to upsert", func(t *testing.T) {
		row := make(types.EventDataRow)
		row["c1"] = "v1"
		row["c2"] = "v2"

		blockData := sqlsol.NewBlockData()
		blockData.AddRow("TEST_TABLE", row)
		blockData.SetBlockID("99")

		hasRows := blockData.PendingRows("99")

		require.Equal(t, true, hasRows)
	})

	t.Run("successfully returns false if a given block does not have pending rows to upsert", func(t *testing.T) {
		row := make(types.EventDataRow)
		row["c1"] = "v1"
		row["c2"] = "v2"

		blockData := sqlsol.NewBlockData()
		blockData.AddRow("TEST_TABLE", row)
		blockData.SetBlockID("99")

		hasRows := blockData.PendingRows("88")

		require.Equal(t, false, hasRows)
	})

	t.Run("successfully returns false if a given block does not exists", func(t *testing.T) {
		blockData := sqlsol.NewBlockData()
		hasRows := blockData.PendingRows("999")

		require.Equal(t, false, hasRows)
	})
}
