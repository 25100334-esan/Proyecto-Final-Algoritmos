package main

import (
	"fmt"
	"log"

	"bplustree-proyecto/api"
	"bplustree-proyecto/database"
)

func main() {
	const rutaBD = "data/dataset.db"
	const puerto = ":8080"

	// 1. Conectar al disco (la conexión queda lista para cuando Vue pida cargar datos).
	db, err := database.Conectar(rutaBD)
	if err != nil {
		log.Fatalf("error al conectar BD: %v", err)
	}
	defer db.Cerrar()

	// 2. Levantar API y esperar que Vue.js llame a /api/inicializar.
	servidor := api.NuevoServidor(db)
	fmt.Printf("Servidor API en http://localhost%s\n", puerto)
	fmt.Println("Esperando configuración desde Vue.js (POST /api/inicializar)...")
	if err := servidor.Iniciar(puerto); err != nil {
		log.Fatalf("error al iniciar servidor: %v", err)
	}
}
