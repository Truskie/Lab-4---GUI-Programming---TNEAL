package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)
	handler.ServeHTTP(rr, req)

	fmt.Println("Home returned: ", rr.Body.String(), rr.Code)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got %v, expected %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	expectedStrings := []string{
		"My name is Terryl Neal",
		"My semester project is a Medical Appointment Scheduling System",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("expected to find %q in response body", expected)
		}
	}
}

func TestAboutHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/about", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(About)
	handler.ServeHTTP(rr, req)

	fmt.Println("About returned: ", rr.Body.String(), rr.Code)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got %v, expected %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	expectedStrings := []string{
		"About Page",
		"My name is Terryl Neal",
		"Bachelor's in IT at UB",
		"catch up in some respects",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("expected to find %q in response", expected)
		}
	}
}

func TestContactHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/contact", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Contact)
	handler.ServeHTTP(rr, req)

	//Debug to see stuff
	fmt.Println("Contact returned: ", rr.Body.String(), rr.Code)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got %v, expected %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	expectedStrings := []string{
		"Contact Page",
		"2024161285@ub.edu.bz",
		"614-1850",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("expected to find %q in contact response", expected)
		}
	}
}

func TestHobbyHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/hobby", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Hobby)
	handler.ServeHTTP(rr, req)

	//Debug to see stuff
	fmt.Println("Hobby returned: ", rr.Body.String(), rr.Code)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got %v, expected %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	expectedStrings := []string{
		"Hobby Page",
		"fishing",
		"do it at 5",
		"destressor",
		"appreciation for the world",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("expected to find %q in hobby response", expected)
		}
	}
}
