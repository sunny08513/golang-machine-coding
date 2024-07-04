package main

import (
	"errors"
	"fmt"
)

type Document struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Owner   string `json:"owner"`
}

type Service struct {
	Documents   map[string]*Document         `json:"document"`
	Permissions map[string]map[string]string `json:"permission"` //document -> (user -> permission)
}

func NewService() *Service {
	docs := make(map[string]*Document)
	permissions := make(map[string]map[string]string)

	return &Service{
		Documents:   docs,
		Permissions: permissions,
	}
}

func (s *Service) CreateDocument(name, content, owner string) error {
	if _, ok := s.Documents[name]; ok {
		return errors.New("document already exist with same name")
	}

	doc := &Document{
		Name:    name,
		Content: content,
		Owner:   owner,
	}
	s.Documents[name] = doc
	s.Permissions[name] = map[string]string{owner: "owner"}
	return nil
}

func (s *Service) ReadDocument(name, user string) (string, error) {
	if doc, exist := s.Documents[name]; exist {
		if per, hasAccess := s.Permissions[name][user]; hasAccess {
			if per == "read" || per == "write" || per == "owner" {
				return doc.Content, nil
			} else {
				return "", errors.New("access denied")
			}
		}
	}
	return "", errors.New("document not found")
}

// EditDocument edits a document if the user has edit access.
func (s *Service) EditDocument(user, name, newContent string) error {
	if doc, exists := s.Documents[name]; exists {
		if perm, hasAccess := s.Permissions[name][user]; hasAccess && (perm == "edit" || perm == "owner") {
			doc.Content = newContent
			return nil
		}
		return errors.New("access denied")
	}
	return errors.New("document not found")
}

// DeleteDocument deletes a document if the user is the owner.
func (s *Service) DeleteDocument(user, name string) error {
	if doc, exists := s.Documents[name]; exists {
		if doc.Owner == user {
			delete(s.Documents, name)
			delete(s.Permissions, name)
			return nil
		}
		return errors.New("access denied")
	}
	return errors.New("document not found")
}

// GrantAccess grants read or edit access to a user for a document.
func (s *Service) GrantAccess(owner, name, user, access string) error {
	if doc, exists := s.Documents[name]; exists {
		if doc.Owner == owner {
			if access != "read" && access != "edit" {
				return errors.New("invalid access type")
			}
			s.Permissions[name][user] = access
			return nil
		}
		return errors.New("access denied")
	}
	return errors.New("document not found")
}

func main() {
	service := NewService()

	// Creating documents
	if err := service.CreateDocument("alice", "doc1", "This is the content of doc1"); err != nil {
		fmt.Println(err)
	}
	if err := service.CreateDocument("bob", "doc2", "This is the content of doc2"); err != nil {
		fmt.Println(err)
	}

	// Granting access
	if err := service.GrantAccess("alice", "doc1", "bob", "read"); err != nil {
		fmt.Println(err)
	}
	if err := service.GrantAccess("alice", "doc1", "charlie", "edit"); err != nil {
		fmt.Println(err)
	}

	// Reading documents
	if content, err := service.ReadDocument("bob", "doc1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Bob reads doc1:", content)
	}

	// Editing documents
	if err := service.EditDocument("charlie", "doc1", "New content for doc1"); err != nil {
		fmt.Println(err)
	}
	if content, err := service.ReadDocument("alice", "doc1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Alice reads updated doc1:", content)
	}

	// Deleting documents
	if err := service.DeleteDocument("bob", "doc1"); err != nil {
		fmt.Println(err)
	}
	if err := service.DeleteDocument("alice", "doc1"); err != nil {
		fmt.Println(err)
	}

	// Checking if the document is deleted
	if content, err := service.ReadDocument("alice", "doc1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Alice reads doc1:", content)
	}
}
