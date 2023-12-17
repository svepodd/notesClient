package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"Notes/models/dto"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Println("Please enter the action you want to perform:\n1. Add new note\n2. Get note by id\n3. Update note by id\n4. Delete note by id\n5. Get all notes\n6. Exit")
		var action string
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		action, _ = reader.ReadString('\n')
		action = strings.TrimSpace(action)
		switch action {
		case "1":
			noteAdd()
		case "2":
			NoteGet()
		case "3":
			noteUpdate()
		case "4":
			noteDelete()
		case "5":
			notesGetAll()
		case "6":
			return
		default:
			fmt.Println("Wrong action")
		}
	}
}

func noteAdd() {
	note := dto.NewNote()
	fmt.Println("Please enter the note data:")

	for note.Name == "" {
		fmt.Println("Name:")
		reader := bufio.NewReader(os.Stdin)
		note.Name, _ = reader.ReadString('\n')
		note.Name = strings.TrimSpace(note.Name)
	}

	for note.LastName == "" {
		fmt.Println("Last name:")
		reader := bufio.NewReader(os.Stdin)
		note.LastName, _ = reader.ReadString('\n')
		note.LastName = strings.TrimSpace(note.LastName)
	}

	for note.Note == "" {
		fmt.Println("Note:")
		reader := bufio.NewReader(os.Stdin)
		note.Note, _ = reader.ReadString('\n')
		note.Note = strings.TrimSpace(note.Note)
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error json.Marshal():", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func NoteGet() {
	note := dto.NewNote()

	fmt.Println("Enter ID of a note you want to get:")
	fmt.Scanln(&note.ID)

	if note.ID < 1 {
		fmt.Println("Error: ID Must be valid")
		return
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func noteUpdate() {
	note := dto.NewNote()

	for note.ID < 1 {
		fmt.Println("Note ID:")
		fmt.Scanln(&note.ID)
	}

	fmt.Println("Please enter updated data:")

	fmt.Println("Name:")
	reader := bufio.NewReader(os.Stdin)
	note.Name, _ = reader.ReadString('\n')
	note.Name = strings.TrimSpace(note.Name)

	fmt.Println("Last name:")
	reader = bufio.NewReader(os.Stdin)
	note.LastName, _ = reader.ReadString('\n')
	note.LastName = strings.TrimSpace(note.LastName)

	fmt.Println("Note:")
	reader = bufio.NewReader(os.Stdin)
	note.Note, _ = reader.ReadString('\n')
	note.Note = strings.TrimSpace(note.Note)

	if note.Name == "" && note.LastName == "" && note.Note == "" {
		fmt.Println("Error: all fields are empty")
		return
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/update", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Response Body:", err)
		return
	}

	ResponseHandler(body)
}

func noteDelete() {
	note := dto.NewNote()
	fmt.Println("Please enter note ID to delete it:")
	fmt.Scanln(&note.ID)

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/delete", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func notesGetAll() {
	resp, err := http.Post("http://localhost:8080/get-all", "application/json", bytes.NewBuffer(nil))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)

}

func ResponseHandler(body []byte) {
	resp := dto.Response{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error in response:", err)
		return
	}

	if resp.Error != "" {
		fmt.Println(resp.Result)
		return
	}

	fmt.Println("")

	if resp.Data != nil {
		data := []dto.Note{}
		err = json.Unmarshal(resp.Data, &data)
		if err != nil {
			data := dto.Note{}
			err = nil
			err = json.Unmarshal(resp.Data, &data)
			if err != nil {
				fmt.Println("Error in response:", err)
				return
			}
			PrintNote(data)
			fmt.Println()
			fmt.Println()
			return
		}

		for _, note := range data {
			PrintNote(note)
		}
	}
	fmt.Println()
	fmt.Println()
}

func PrintNote(note dto.Note) {
	fmt.Println("ID:", note.ID)
	fmt.Println("Creator's Name:", note.Name)
	fmt.Println("Creator's Last name:", note.LastName)
	fmt.Println("Note:", note.Note)
	fmt.Println()
}
