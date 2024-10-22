package utils

func Float64ToFloat32(embed []float64) []float32 {
	var embed32 []float32
	for _, e := range embed {
		embed32 = append(embed32, float32(e))
	}

	return embed32
}
