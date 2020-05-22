package main

func FilterMinPrice(island Island, MinPrice int) bool {
	if island.TurnipPrice < MinPrice {
		return false
	}

	return true
}

func FilterMaxPrice(island Island, MaxPrice int) bool {
	if island.TurnipPrice > MaxPrice {
		return false
	}

	return true
}

func FilterExcludePrices(island Island, ExcludePrices []int) bool {
	for _, price := range ExcludePrices {
		if island.TurnipPrice == price {
			return false
		}
	}

	return true
}

func FilterQueueSize(island Island, MaxInQueue int) bool {
	if island.InQueue > MaxInQueue {
		return false
	}

	return true
}
