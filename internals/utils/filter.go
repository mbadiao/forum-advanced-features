package utils

import (
	"fmt"
	"forum/internals/database"
	"strconv"
	"strings"
)

func CheckCategory(category []string) bool {
	Mapcategory := map[string]bool{
		"All":     true,
		"Tech":    true,
		"Actu":    true,
		"Mode":    true,
		"Sport":   true,
		"Edu":     true,
		"Like":    true,
		"Created": true,
	}
	found := true
	for _, v := range category {
		if !CompareCategorys(Mapcategory, v) {
			found = false
		}
	}
	return found
}

func CompareCategorys(categoriesMap map[string]bool, categoryToCheck string) bool {
	_, found := categoriesMap[categoryToCheck]
	return found
}

func QueryFilter(categorypost, createdlikedpost []string, foundAll bool, Isconnected bool, user database.User) (string, string) {
	var query string
	FoundQuery := CheckQuery(categorypost, createdlikedpost, foundAll)
	switch {
	case FoundQuery == "all":
		query = "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts ORDER BY creation_date DESC"
	case FoundQuery == "category":
		query = "SELECT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date FROM Posts p INNER JOIN PostCategories pc ON p.post_id = pc.post_id INNER JOIN Categories c ON pc.category_id = c.category_id WHERE c.name IN ("
		placeholders := make([]string, len(categorypost))
		for i := range categorypost {
			placeholders[i] = "?"
		}
		query += strings.Join(placeholders, ",") + ") ORDER BY creation_date DESC"
	case FoundQuery == "like" && Isconnected:
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date  FROM Posts p  JOIN LikesDislikes ld ON p.post_id = ld.post_id  WHERE ld.user_id = " + strconv.Itoa(user.UserID) + "  AND ld.liked = TRUE ORDER BY creation_date DESC"
	case FoundQuery == "create" && Isconnected:
		query = "SELECT post_id, user_id, title, PhotoURL, content, creation_date FROM Posts WHERE user_id =" + strconv.Itoa(user.UserID)
	case FoundQuery == "createlike" && Isconnected:
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date  FROM Posts p  LEFT JOIN LikesDislikes ld ON p.post_id = ld.post_id  LEFT JOIN Users u ON p.user_id = u.user_id  WHERE p.user_id = " + strconv.Itoa(user.UserID) + " AND ld.user_id = " + strconv.Itoa(user.UserID) + "  ORDER BY p.creation_date DESC"
	case FoundQuery == "likecategory" && Isconnected:
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date FROM Posts p INNER JOIN PostCategories pc ON p.post_id = pc.post_id INNER JOIN Categories c ON pc.category_id = c.category_id INNER JOIN LikesDislikes ld ON p.post_id = ld.post_id AND ld.liked = TRUE WHERE c.name IN ("
		placeholders := make([]string, len(categorypost))
		for i := range categorypost {
			placeholders[i] = "?"
		}
		query += strings.Join(placeholders, ",") + ")"
		query += "GROUP BY p.post_id ORDER BY p.creation_date DESC"
	case FoundQuery == "createcategory" && Isconnected:
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date FROM Posts p INNER JOIN Users u ON p.user_id = u.user_id INNER JOIN PostCategories pc ON p.post_id = pc.post_id INNER JOIN Categories c ON pc.category_id = c.category_id WHERE u.user_id = " + strconv.Itoa(user.UserID) + " AND c.name IN ("
		placeholders := make([]string, len(categorypost))
		for i := range categorypost {
			placeholders[i] = "?"
		}
		query += strings.Join(placeholders, ",") + ") ORDER BY creation_date DESC"
	case FoundQuery == "createlikecategory" && Isconnected:
		query = "SELECT DISTINCT p.post_id, p.user_id, p.title, p.PhotoURL, p.content, p.creation_date FROM Posts p INNER JOIN Users u ON p.user_id = u.user_id LEFT JOIN PostCategories pc ON p.post_id = pc.post_id LEFT JOIN Categories c ON pc.category_id = c.category_id LEFT JOIN LikesDislikes ld ON p.post_id = ld.post_id AND ld.liked = TRUE WHERE u.user_id = " + strconv.Itoa(user.UserID) + " AND c.name IN ("
		placeholders := make([]string, len(categorypost))
		for i := range categorypost {
			placeholders[i] = "?"
		}
		query += strings.Join(placeholders, ",") + ")"
		query += " GROUP BY p.post_id, u.username ORDER BY p.creation_date DESC"
	default:
		return "", "err"
	}
	return query, ""
}

func SplitFilter(checkedvalue []string) ([]string, []string, bool) {
	var categorypost []string
	var createdlikedpost []string
	foundAll := false
	for _, v := range checkedvalue {
		if v == "All" {
			foundAll = true
		} else if v == "Like" || v == "Created" {
			createdlikedpost = append(createdlikedpost, v)
		} else {
			categorypost = append(categorypost, v)
		}
	}
	return categorypost, createdlikedpost, foundAll
}

func CheckQuery(categorypost, createdlikedpost []string, foundAll bool) string {
	if foundAll {
		fmt.Println("trie sur all")
		return "all"
	} else {
		if len(categorypost) == 0 {
			if len(createdlikedpost) == 2 {
				fmt.Println("trie sur post creer et like")
				return "createlike"
			} else if createdlikedpost[0] == "Like" {
				fmt.Println("trie sur like")
				return "like"
			} else {
				fmt.Println("trie sur creer")
				return "create"
			}
		} else if len(createdlikedpost) == 0 {
			fmt.Println("trie sur post category")
			return "category"
		} else if len(createdlikedpost) == 2 {
			fmt.Println("trie sur post creer et like et category")
			return "createlikecategory"
		} else if createdlikedpost[0] == "Like" {
			fmt.Println("trie sur like et category")
			return "likecategory"
		} else {
			fmt.Println("trie sur create et category")
			return "createcategory"
		}
	}
}

func Isconnected(user database.User) bool {
	if user.UserID == 0 {
		return false
	} else {
		return true
	}
}
