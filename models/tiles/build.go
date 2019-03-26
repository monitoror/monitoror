package tiles

const BuildTileType TileType = "BUILD"

type (
	// Used by Jenkins, Gitlab, Teamcity ... as response structure
	BuildTile struct {
		*Tile

		//TODO
	}
)

func NewBuildTile(subType TileSubType) *BuildTile {
	return &BuildTile{
		Tile: &Tile{
			Type:    BuildTileType,
			SubType: subType,
		},
	}
}
