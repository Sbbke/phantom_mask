package controllers

import(
	"PhantomBE/global"
	"PhantomBE/app/api"
	"strings"
	"context"
	"fmt"
)


// Helper method to search pharmacies
func (pc *PharmacyController) searchPharmacies(ctx context.Context, query string) ([]api.SearchResult, error) {
	var pharmacies []global.Pharmacy
	orderClause := fmt.Sprintf(
        "CASE WHEN name ILIKE '%s' THEN 1 WHEN name ILIKE '%s' THEN 2 ELSE 3 END",
        query,
        query+"%",
    )
	// Use parameterized query to prevent SQL injection
	err := pc.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Order(orderClause).
		Limit(50). // Limit results to prevent memory issues
		Find(&pharmacies).Error
	
	if err != nil {
		return nil, err
	}

	results := make([]api.SearchResult, len(pharmacies))
	for i, pharmacy := range pharmacies {
		relevance := calculateRelevance(query, pharmacy.Name)
		results[i] = api.SearchResult{
			Type:      "pharmacy",
			ID:        pharmacy.ID,
			Name:      pharmacy.Name,
			Relevance: relevance,
		}
	}
	
	return results, nil
}

// Helper method to search users
func (pc *PharmacyController) searchUsers(ctx context.Context, query string) ([]api.SearchResult, error) {
	var users []global.User
	orderClause := fmt.Sprintf(
        "CASE WHEN name ILIKE '%s' THEN 1 WHEN name ILIKE '%s' THEN 2 ELSE 3 END",
        query,
        query+"%",
    )
	err := pc.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Order(orderClause).
		Limit(50).
		Find(&users).Error
	
	if err != nil {
		return nil, err
	}

	results := make([]api.SearchResult, len(users))
	for i, user := range users {
		relevance := calculateRelevance(query, user.Name)
		results[i] = api.SearchResult{
			Type:      "user",
			ID:        user.ID,
			Name:      user.Name,
			Relevance: relevance,
		}
	}
	
	return results, nil
}

// Helper method to search masks
func (pc *PharmacyController) searchMasks(ctx context.Context, query string) ([]api.SearchResult, error) {
	var masks []global.Mask
	orderClause := fmt.Sprintf(
        "CASE WHEN name ILIKE '%s' THEN 1 WHEN name ILIKE '%s' THEN 2 ELSE 3 END",
        query,
        query+"%",
    )
	err := pc.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Order(orderClause).
		Limit(50).
		Find(&masks).Error
	
	if err != nil {
		return nil, err
	}

	results := make([]api.SearchResult, len(masks))
	for i, mask := range masks {
		relevance := calculateRelevance(query, mask.Name)
		results[i] = api.SearchResult{
			Type:       "mask",
			ID:         mask.ID,
			Name:       mask.Name,
			Price:      &mask.Price,
			PharmacyID: &mask.PharmacyID,
			Relevance:  relevance,
		}
	}
	
	return results, nil
}

// Helper function to calculate relevance score based on search query
func calculateRelevance(query, target string) float64 {
	queryLower := strings.ToLower(strings.TrimSpace(query))
	targetLower := strings.ToLower(strings.TrimSpace(target))
	
	if queryLower == "" || targetLower == "" {
		return 0.0
	}
	
	// Exact match gets highest score
	if queryLower == targetLower {
		return 100.0
	}
	
	// Starts with query gets high score
	if strings.HasPrefix(targetLower, queryLower) {
		return 90.0
	}
	
	// Ends with query gets good score
	if strings.HasSuffix(targetLower, queryLower) {
		return 80.0
	}
	
	// Contains query as whole word gets medium-high score
	if strings.Contains(" "+targetLower+" ", " "+queryLower+" ") {
		return 70.0
	}
	
	// Contains query gets medium score
	if strings.Contains(targetLower, queryLower) {
		return 60.0
	}
	
	// Calculate similarity based on common characters
	commonChars := 0
	for _, char := range queryLower {
		if strings.ContainsRune(targetLower, char) {
			commonChars++
		}
	}
	
	if commonChars > 0 {
		return float64(commonChars) / float64(len(queryLower)) * 40.0
	}
	
	return 0.0
}



