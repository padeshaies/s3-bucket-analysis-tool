package helpers

import (
	"fmt"
	"math"
)

func CalculateObjectsCostByStorageType(storageType, region string, sizeInBytes, objectNumber int) (float64, error) {
	switch storageType {
	case "STANDARD":
		return calculateStandardCost(sizeInBytes, region)
	case "INTELLIGENT_TIERING":
		// TODO: Implement the cost calculation for the INTELLIGENT_TIERING storage type
		return 0.0, nil
	case "REDUCED_REDUNDANCY":
		return calculateReducedRedundancyCost(sizeInBytes, region)
	case "GLACIER", "GLACIER_IR", "DEEP_ARCHIVE", "STANDARD_IA", "ONEZONE_IA", "EXPRESS_ONEZONE":
		return calculateCost(sizeInBytes, region, storageType)
	case "OUTPOSTS":
		// Is it 0$ because it's an on-premises storage and has been paid upfront?
		return 0.0, nil
	case "SNOW":
		// TODO: Implement the cost calculation for the SNOW storage type
		return 0.0, nil
	}

	return 0.0, fmt.Errorf("invalid storage type")
}

func calculateStandardCost(sizeInBytes int, region string) (float64, error) {
	totalCost := 0.0

	// First 50 GB
	multiplier, err := getCostMuliplier(region, "STANDARD_<50GB")
	if err != nil {
		return 0.0, err
	}
	totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, 50))

	// Next 450 GB
	if sizeInBytes > gbToBytes(50) {
		multiplier, err := getCostMuliplier(region, "STANDARD_<450GB")
		if err != nil {
			return 0.0, err
		}
		totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, 500)-50)
	}

	// Over 500 GB
	if sizeInBytes > gbToBytes(500) {
		multiplier, err := getCostMuliplier(region, "STANDARD_<500GB")
		if err != nil {
			return 0.0, err
		}
		totalCost += multiplier * float64(sizeInBytes/1024/1024/1024-500)
	}

	return round2(totalCost), nil
}

func calculateReducedRedundancyCost(sizeInBytes int, region string) (float64, error) {
	totalCost := 0.0

	// First 1 TB
	multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_<1TB")
	if err != nil {
		return 0.0, err
	}
	totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, tbToGb(1)))

	// Next 49 TB
	if sizeInBytes > tbToBytes(1) {
		multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_<49TB")
		if err != nil {
			return 0.0, err
		}

		totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, tbToGb(50))-tbToGb(1))
	}

	// Next 450 TB
	if sizeInBytes > tbToBytes(50) {
		multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_<450TB")
		if err != nil {
			return 0.0, err
		}

		totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, tbToGb(500))-tbToGb(50))
	}

	// Next 500 TB
	if sizeInBytes > tbToBytes(500) {
		multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_<1000TB")
		if err != nil {
			return 0.0, err
		}

		totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, tbToGb(1000))-tbToGb(500))
	}

	// Next 4000 TB
	if sizeInBytes > tbToBytes(1000) {
		multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_<5000TB")
		if err != nil {
			return 0.0, err
		}
		totalCost += multiplier * float64(min(sizeInBytes/1024/1024/1024, tbToGb(5000))-tbToGb(1000))
	}

	// Over 5000 TB
	if sizeInBytes > tbToBytes(5000) {
		multiplier, err := getCostMuliplier(region, "REDUCED_REDUNDANCY_>5000TB")
		if err != nil {
			return 0.0, err
		}
		totalCost += multiplier * float64(sizeInBytes/1024/1024/1024-tbToGb(5000))
	}

	return round2(totalCost), nil
}

func calculateCost(sizeInBytes int, region, storageType string) (float64, error) {
	multiplier, err := getCostMuliplier(region, storageType)
	if err != nil {
		return 0.0, err
	}

	return round2(multiplier * float64(sizeInBytes/1024/1024/1024)), nil
}

func gbToBytes(i int) int {
	return i * 1024 * 1024 * 1024
}

func tbToGb(i int) int {
	return i * 1024
}

func tbToBytes(i int) int {
	return i * 1024 * 1024 * 1024 * 1024
}

func round2(n float64) float64 {
	return math.Round(n*100) / 100
}

func getCostMuliplier(region, multiplier string) (float64, error) {
	switch region {
	case "us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-north-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.023, nil
		case "STANDARD_<450GB":
			return 0.022, nil
		case "STANDARD_<500GB":
			return 0.021, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.024, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0236, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0232, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0228, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0224, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.022, nil
		case "STANDARD_IA":
			return 0.0125, nil
		case "GLACIER":
			return 0.0036, nil
		case "GLACIER_IR":
			return 0.004, nil
		case "DEEP_ARCHIVE":
			return 0.00099, nil
		case "EXPRESS_ONEZONE":
			return 0.016, nil
		case "ONEZONE_IA":
			return 0.01, nil
		default:
			return 0.0, fmt.Errorf("invalid multiplier")
		}
	case "ca-central-1", "ca-west-1", "il-central-1", "me-south-1", "me-central-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.025, nil
		case "STANDARD_<450GB":
			return 0.024, nil
		case "STANDARD_<500GB":
			return 0.023, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0264, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.026, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0255, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0251, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0246, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0242, nil
		case "STANDARD_IA":
			return 0.0138, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.01104, nil
		default:
			return 0.0, fmt.Errorf("invalid multiplier")
		}
	case "mx-central-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.02415, nil
		case "STANDARD_<450GB":
			return 0.0231, nil
		case "STANDARD_<500GB":
			return 0.02205, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0252, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.02478, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.02436, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.02394, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.02352, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0231, nil
		case "STANDARD_IA":
			return 0.013125, nil
		case "GLACIER":
			return 0.00378, nil
		case "GLACIER_IR":
			return 0.0042, nil
		case "DEEP_ARCHIVE":
			return 0.002, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.0105, nil
		}
	case "us-gov-east-1", "us-gov-west-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.039, nil
		case "STANDARD_<450GB":
			return 0.037, nil
		case "STANDARD_<500GB":
			return 0.0355, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0312, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0306, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0301, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0296, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0291, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0285, nil
		case "STANDARD_IA":
			return 0.02, nil
		case "GLACIER":
			return 0.0054, nil
		case "GLACIER_IR":
			return 0.0064, nil
		case "DEEP_ARCHIVE":
			return 0.0024, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.016, nil
		}
	case "af-south-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.0274, nil
		case "STANDARD_<450GB":
			return 0.0262, nil
		case "STANDARD_<500GB":
			return 0.025, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.024, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0236, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0232, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0228, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0224, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.022, nil
		case "STANDARD_IA":
			return 0.0149, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.0119, nil
		}
	case "ap-east-1", "ap-south-2", "ap-southeast-3", "ap-southeast-4", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.025, nil
		case "STANDARD_<450GB":
			return 0.024, nil
		case "STANDARD_<500GB":
			return 0.023, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.024, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0236, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0232, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0228, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0224, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.022, nil
		case "STANDARD_IA":
			return 0.0138, nil
		case "GLACIER":
			return 0.0045, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.002, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.011, nil
		}
	case "ap-south-1", "ap-northeast-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.025, nil
		case "STANDARD_<450GB":
			return 0.024, nil
		case "STANDARD_<500GB":
			return 0.023, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.024, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0236, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0232, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0228, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0224, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.022, nil
		case "STANDARD_IA":
			return 0.0138, nil
		case "GLACIER":
			return 0.0045, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.002, nil
		case "EXPRESS_ONEZONE":
			return 0.18, nil
		case "ONEZONE_IA":
			return 0.011, nil
		}
	case "ap-southeast-5", "ap-southeast-7":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.0225, nil
		case "STANDARD_<450GB":
			return 0.0216, nil
		case "STANDARD_<500GB":
			return 0.0207, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0216, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.02124, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.02088, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.02052, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.02016, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0198, nil
		case "STANDARD_IA":
			return 0.01242, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.0045, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.0099, nil
		}
	case "eu-central-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.0245, nil
		case "STANDARD_<450GB":
			return 0.0235, nil
		case "STANDARD_<500GB":
			return 0.0225, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.026, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0255, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0251, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0247, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0242, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0238, nil
		case "STANDARD_IA":
			return 0.0135, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.0108, nil
		}
	case "eu-west-2", "eu-south-1", "eu-west-3":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.024, nil
		case "STANDARD_<450GB":
			return 0.023, nil
		case "STANDARD_<500GB":
			return 0.022, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0252, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0248, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0244, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0239, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0235, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0231, nil
		case "STANDARD_IA":
			return 0.0131, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.01048, nil
		}
	case "eu-south-2":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.023, nil
		case "STANDARD_<450GB":
			return 0.022, nil
		case "STANDARD_<500GB":
			return 0.021, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.024, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.0236, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0232, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0228, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0224, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.022, nil
		case "STANDARD_IA":
			return 0.0125, nil
		case "GLACIER":
			return 0.00405, nil
		case "GLACIER_IR":
			return 0.005, nil
		case "DEEP_ARCHIVE":
			return 0.0018, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.01, nil
		}
	case "eu-central-2":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.02695, nil
		case "STANDARD_<450GB":
			return 0.02585, nil
		case "STANDARD_<500GB":
			return 0.02475, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0286, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.02805, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.02761, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.02717, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.02662, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.02618, nil
		case "STANDARD_IA":
			return 0.01485, nil
		case "GLACIER":
			return 0.004455, nil
		case "GLACIER_IR":
			return 0.0055, nil
		case "DEEP_ARCHIVE":
			return 0.00198, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.01188, nil
		}
	case "sa-east-1":
		switch multiplier {
		case "STANDARD_<50GB":
			return 0.0405, nil
		case "STANDARD_<450GB":
			return 0.039, nil
		case "STANDARD_<500GB":
			return 0.037, nil
		case "REDUCED_REDUNDANCY_<1TB":
			return 0.0326, nil
		case "REDUCED_REDUNDANCY_<49TB":
			return 0.032, nil
		case "REDUCED_REDUNDANCY_<450TB":
			return 0.0315, nil
		case "REDUCED_REDUNDANCY_<1000TB":
			return 0.0309, nil
		case "REDUCED_REDUNDANCY_<5000TB":
			return 0.0304, nil
		case "REDUCED_REDUNDANCY_>5000TB":
			return 0.0299, nil
		case "STANDARD_IA":
			return 0.0221, nil
		case "GLACIER":
			return 0.00765, nil
		case "GLACIER_IR":
			return 0.0083, nil
		case "DEEP_ARCHIVE":
			return 0.0032, nil
		case "EXPRESS_ONEZONE":
			return 0.0, fmt.Errorf("express onezone is not available in this region")
		case "ONEZONE_IA":
			return 0.0177, nil
		}
	}

	// Should return an error if the region is not supported, but for the sake of the example, we will return 0.0
	// return 0.0, fmt.Errorf("invalid region (received: %v)", region)
	return 0.0, nil
}
