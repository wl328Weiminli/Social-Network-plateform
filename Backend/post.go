package main

import (
    "fmt"
    "mime/multipart"
    "reflect"

    "github.com/olivere/elastic"
    "github.com/pborman/uuid"
)

const (
    POST_INDEX  = "post"
)


type Post struct {
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`
    Type    string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
    query := elastic.NewTermQuery("user", user)
    searchResult, err := readFromES(query, POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND")
    if keywords == "" {
        query.ZeroTermsQuery("all")
    }
    searchResult, err := readFromES(query, POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
    var ptype Post
    var posts []Post

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
        if p, ok := item.(Post); ok {
            posts = append(posts, p)
        }
    }
    return posts
}

func savePost(post *Post, file multipart.File) error {
    id := uuid.New()
    medialink, err := saveToGCS(file, id)
    if err != nil {
        return err
    }
    post.Url = medialink

    err = saveToES(post, POST_INDEX, id)
    if err != nil {
        return err
    }
    fmt.Printf("Post is saved to Elasticsearch: %s\n", post.Message)
    return nil
}




