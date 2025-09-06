package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Message.GetConversation() != "" {
			fmt.Printf("📱 Nachricht von %s: %s\n", v.Info.Sender, v.Message.GetConversation())
		}
	case *events.Receipt:
		if v.Type == types.ReceiptTypeRead || v.Type == types.ReceiptTypeReadSelf {
			fmt.Printf("✓ Nachricht gelesen von %s\n", v.SourceString())
		}
	}
}

// Funktion zum Senden einer einzelnen Nachricht (für CLI-Modus)
func sendSingleMessage(ctx context.Context, client *whatsmeow.Client, recipient, message string) error {
	// Nummer/Gruppen-ID formatieren
	if !strings.Contains(recipient, "@") {
		if strings.Contains(recipient, "-") && len(recipient) > 15 {
			// Sieht aus wie eine Gruppen-ID
			recipient = recipient + "@g.us"
		} else {
			// Normale Telefonnummer
			recipient = strings.TrimPrefix(recipient, "+")
			recipient = recipient + "@s.whatsapp.net"
		}
	}

	// JID erstellen
	jid, err := types.ParseJID(recipient)
	if err != nil {
		return fmt.Errorf("ungültige Nummer/Gruppen-ID: %v", err)
	}

	// Nachricht senden
	_, err = client.SendMessage(ctx, jid, &waE2E.Message{
		Conversation: &message,
	})

	if err != nil {
		return fmt.Errorf("fehler beim Senden: %v", err)
	}

	if strings.Contains(jid.Server, "g.us") {
		fmt.Printf("✅ Gruppennachricht gesendet: %s\n", message)
	} else {
		fmt.Printf("✅ Nachricht an %s gesendet: %s\n", jid.User, message)
	}

	return nil
}

// Optional: Debug-Modus aus Environment-Variable
func getLogLevel() string {
	if os.Getenv("WHATSAPP_DEBUG") == "1" {
		return "DEBUG"
	}
	return "WARN" // Weniger verbose für normale Nutzung
}

func main() {
	// Context für die Anwendung
	ctx := context.Background()

	// CLI-Argumente prüfen
	args := os.Args[1:]

	// Hilfe anzeigen
	if len(args) == 1 && (args[0] == "help" || args[0] == "-h" || args[0] == "--help") {
		fmt.Println("🚀 WhatsApp Console Client")
		fmt.Println()
		fmt.Println("📋 Verwendung:")
		fmt.Println("  Interaktiv:  ./whatsapp-console")
		fmt.Println("  CLI-Modus:   ./whatsapp-console <befehl> [argumente]")
		fmt.Println()
		fmt.Println("📤 Beispiele:")
		fmt.Println("  Einzelperson: ./whatsapp-console send 491234567890 'Hallo Welt!'")
		fmt.Println("  Gruppe:       ./whatsapp-console send 120363XX...XX@g.us 'Hallo Gruppe!'")
		fmt.Println("  Kontakte:     ./whatsapp-console list")
		fmt.Println("  Alle Kontakte: ./whatsapp-console list all")
		fmt.Println("  Gruppen:      ./whatsapp-console list groups")
		fmt.Println()
		fmt.Println("💡 Hinweise:")
		fmt.Println("  • Beim ersten Start: QR-Code scannen erforderlich")
		fmt.Println("  • Nachrichten mit Leerzeichen in Anführungszeichen setzen")
		fmt.Println("  • Gruppen-IDs: Format 120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us")
		os.Exit(0)
	}

	// Version anzeigen
	if len(args) == 1 && (args[0] == "version" || args[0] == "-v" || args[0] == "--version") {
		fmt.Printf("WhatsApp Console Client %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		os.Exit(0)
	}

	// CLI-Variablen vorbereiten
	var recipient, message string
	var isCliMode bool = len(args) >= 1
	var cliCommand string

	if len(args) >= 1 {
		cliCommand = strings.ToLower(args[0])

		if cliCommand == "send" && len(args) >= 3 {
			// CLI-Modus: ./whatsapp-console send <nummer> <nachricht>
			recipient = args[1]
			message = strings.Join(args[2:], " ") // Alle weiteren Argumente als Nachricht
		}
	}

	// Database für Session-Speicherung
	dbLog := waLog.Stdout("Database", getLogLevel(), true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Printf("❌ Fehler beim Erstellen der Datenbank: %v\n", err)
		return
	}

	// Database-Berechtigungen sichern (nur Owner kann lesen/schreiben)
	if err := os.Chmod("whatsapp.db", 0600); err != nil {
		fmt.Printf("⚠️ Warnung: Konnte Datenbankberechtigungen nicht setzen: %v\n", err)
	}

	// Device Store
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		fmt.Printf("❌ Fehler beim Laden des Geräts: %v\n", err)
		return
	}

	// Client Logger
	clientLog := waLog.Stdout("Client", getLogLevel(), true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Signal Handler für sauberes Beenden
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n👋 Beende Programm...")

		// Graceful disconnect mit Timeout
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client.Disconnect()

		// Warten auf sauberes Herunterfahren
		select {
		case <-disconnectCtx.Done():
			fmt.Println("⏰ Timeout beim Herunterfahren")
		case <-time.After(1 * time.Second):
			fmt.Println("✅ Sauber beendet")
		}

		os.Exit(0)
	}()

	// QR Code handling wenn nicht eingeloggt
	if client.Store.ID == nil {
		if len(args) >= 3 && args[0] == "send" {
			fmt.Println("❌ Erstes Login erforderlich!")
			fmt.Println("💡 Führen Sie zuerst './whatsapp-console' aus um QR-Code zu scannen")
			os.Exit(1)
		}

		fmt.Println("📱 Erstes Login erforderlich!")
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			fmt.Printf("❌ Verbindungsfehler: %v\n", err)
			return
		}

		fmt.Println("📷 Scannen Sie den QR-Code mit WhatsApp:")
		fmt.Println()
		for evt := range qrChan {
			if evt.Event == "code" {
				// QR-Code im Terminal anzeigen
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println()
				fmt.Println("💡 QR-Code wurde auch als 'qr.png' gespeichert")

				// QR-Code auch als PNG-Datei speichern (optional)
				// Sie können diesen Code auskommentieren wenn gewünscht
				/*
					err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr.png")
					if err == nil {
						fmt.Println("📁 QR-Code gespeichert als qr.png")
					}
				*/

				fmt.Println("🔗 QR-Code Text:", evt.Code)
			} else {
				fmt.Printf("QR Event: %s\n", evt.Event)
			}
		}
	} else {
		// Bereits eingeloggt, nur verbinden
		if len(args) < 3 || args[0] != "send" {
			fmt.Println("🔗 Verbinde mit WhatsApp...")
		}
		err = client.Connect()
		if err != nil {
			fmt.Printf("❌ Verbindungsfehler: %v\n", err)
			return
		}
	}

	// Warten bis verbunden
	if !isCliMode {
		fmt.Println("⏳ Warte auf Verbindung...")
	}
	if !client.WaitForConnection(10 * time.Second) {
		fmt.Println("❌ Timeout beim Verbinden")
		return
	}

	if !isCliMode {
		fmt.Println("✅ Erfolgreich mit WhatsApp verbunden!")
	}

	// CLI-Modus behandeln
	if isCliMode {
		switch cliCommand {
		case "send":
			if len(args) >= 3 {
				err = sendSingleMessage(ctx, client, recipient, message)
				client.Disconnect()
				if err != nil {
					fmt.Printf("❌ %v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			} else {
				fmt.Println("❌ Verwendung: ./whatsapp-console send <nummer/gruppe> <nachricht>")
				os.Exit(1)
			}

		case "list":
			// CLI list Befehl
			showAll := len(args) > 1 && strings.ToLower(args[1]) == "all"
			showGroups := len(args) > 1 && strings.ToLower(args[1]) == "groups"

			if showGroups {
				// Lade Gruppen
				groups, err := client.GetJoinedGroups()
				if err != nil {
					fmt.Printf("❌ Fehler beim Laden der Gruppen: %v\n", err)
					os.Exit(1)
				}

				if len(groups) == 0 {
					fmt.Println("📭 Keine Gruppen gefunden.")
					os.Exit(0)
				}

				fmt.Printf("🏘️ Gefundene Gruppen (%d):\n\n", len(groups))
				for i, group := range groups {
					fmt.Printf("%d. %s\n", i+1, group.Name)
					fmt.Printf("   � ID: %s\n", group.JID.String())
					if group.Topic != "" {
						fmt.Printf("   📄 Beschreibung: %s\n", group.Topic)
					}
					fmt.Printf("   👥 Teilnehmer: %d\n\n", len(group.Participants))
				}
			} else {
				// Kontakte auflisten
				contacts, err := client.Store.Contacts.GetAllContacts(ctx)
				if err != nil {
					fmt.Printf("❌ Fehler beim Laden der Kontakte: %v\n", err)
					os.Exit(1)
				}

				if len(contacts) == 0 {
					fmt.Println("📭 Keine Kontakte gefunden")
					os.Exit(0)
				}

				if showAll {
					fmt.Printf("👥 Alle %d Kontakte:\n\n", len(contacts))
					for jid, contact := range contacts {
						name := contact.PushName
						if name == "" {
							name = contact.BusinessName
						}
						if name == "" {
							name = "Unbekannt"
						}
						fmt.Printf("📞 %s (%s)\n", name, jid.User)
					}
				} else {
					fmt.Printf("👥 %d Kontakte gefunden:\n\n", len(contacts))
					count := 0
					for jid, contact := range contacts {
						if count >= 20 {
							fmt.Printf("... und %d weitere Kontakte\n", len(contacts)-20)
							fmt.Println("💡 Verwenden Sie './whatsapp-console list all' für alle Kontakte")
							break
						}
						name := contact.PushName
						if name == "" {
							name = contact.BusinessName
						}
						if name == "" {
							name = "Unbekannt"
						}
						fmt.Printf("📞 %s (%s)\n", name, jid.User)
						count++
					}
				}
			}

			client.Disconnect()
			os.Exit(0)

		default:
			fmt.Printf("❌ Unbekannter Befehl: %s\n", cliCommand)
			fmt.Println("💡 Verwenden Sie './whatsapp-console help' für Hilfe")
			os.Exit(1)
		}
	}

	// Interaktiver Modus
	fmt.Println("\n📋 Verfügbare Befehle:")
	fmt.Println("  📤 send <nummer> <nachricht>  - Nachricht senden")
	fmt.Println("  📋 list                       - Chats auflisten")
	fmt.Println("  🚪 quit                       - Programm beenden")
	fmt.Println("\n💡 Beispiel: send 491234567890 Hallo von der Konsole!")
	fmt.Println("\n🔥 CLI-Modus: ./whatsapp-console send 491234567890 'Hallo Welt!'")
	fmt.Println()

	// Konsolen-Input Loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("whatsapp> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 3)
		command := strings.ToLower(parts[0])

		switch command {
		case "quit", "exit", "q":
			fmt.Println("👋 Auf Wiedersehen!")
			client.Disconnect()
			return

		case "send", "s":
			if len(parts) < 3 {
				fmt.Println("❌ Verwendung: send <nummer/gruppen-id> <nachricht>")
				fmt.Println("💡 Einzelperson: send 491234567890 Hallo Welt!")
				fmt.Println("💡 Gruppe: send 120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us Hallo Gruppe!")
				continue
			}

			number := parts[1]
			message := parts[2]

			// Nummer/Gruppen-ID formatieren
			if !strings.Contains(number, "@") {
				if strings.Contains(number, "-") && len(number) > 15 {
					// Sieht aus wie eine Gruppen-ID
					number = number + "@g.us"
					fmt.Printf("🏘️ Sende an Gruppe: %s\n", number)
				} else {
					// Normale Telefonnummer
					if strings.HasPrefix(number, "+") {
						number = number[1:] // '+' entfernen
					}
					number = number + "@s.whatsapp.net"
					fmt.Printf("👤 Sende an Einzelperson: %s\n", number)
				}
			}

			// JID erstellen
			recipient, err := types.ParseJID(number)
			if err != nil {
				fmt.Printf("❌ Ungültige Nummer/Gruppen-ID: %v\n", err)
				continue
			}

			// Nachricht senden
			_, err = client.SendMessage(ctx, recipient, &waE2E.Message{
				Conversation: &message,
			})

			if err != nil {
				fmt.Printf("❌ Fehler beim Senden: %v\n", err)
			} else {
				if strings.Contains(recipient.Server, "g.us") {
					fmt.Printf("✅ Gruppennachricht gesendet: %s\n", message)
				} else {
					fmt.Printf("✅ Nachricht an %s gesendet: %s\n", recipient.User, message)
				}
			}

		case "list", "l":
			// Prüfen ob "list all" eingegeben wurde
			showAll := len(parts) > 1 && strings.ToLower(parts[1]) == "all"
			showGroups := len(parts) > 1 && strings.ToLower(parts[1]) == "groups"

			if showGroups {
				fmt.Println("🏘️ Lade Gruppen-Liste...")
				// Lade Gruppen
				groups, err := client.GetJoinedGroups()
				if err != nil {
					fmt.Printf("❌ Fehler beim Laden der Gruppen: %v\n", err)
					continue
				}

				if len(groups) == 0 {
					fmt.Println("📭 Keine Gruppen gefunden.")
					fmt.Println()
					continue
				}

				fmt.Printf("🏘️ Gefundene Gruppen (%d):\n\n", len(groups))
				for i, group := range groups {
					fmt.Printf("%d. %s\n", i+1, group.Name)
					fmt.Printf("   � ID: %s\n", group.JID.String())
					if group.Topic != "" {
						fmt.Printf("   📄 Beschreibung: %s\n", group.Topic)
					}
					fmt.Printf("   👥 Teilnehmer: %d\n\n", len(group.Participants))
				}
				continue
			}

			if showAll {
				fmt.Println("📋 Lade vollständige Chat-Liste...")
			} else {
				fmt.Println("📋 Lade Chat-Liste...")
			}

			// Alle Kontakte aus dem Store abrufen
			contacts, err := client.Store.Contacts.GetAllContacts(ctx)
			if err != nil {
				fmt.Printf("❌ Fehler beim Laden der Kontakte: %v\n", err)
				continue
			}

			if len(contacts) == 0 {
				fmt.Println("📭 Keine Kontakte gefunden")
				continue
			}

			if showAll {
				fmt.Printf("👥 Alle %d Kontakte:\n\n", len(contacts))

				for jid, contact := range contacts {
					name := contact.PushName
					if name == "" {
						name = contact.BusinessName
					}
					if name == "" {
						name = "Unbekannt"
					}

					fmt.Printf("� %s (%s)\n", name, jid.User)
				}
			} else {
				fmt.Printf("�👥 %d Kontakte gefunden:\n\n", len(contacts))

				count := 0
				for jid, contact := range contacts {
					if count >= 20 { // Nur die ersten 20 anzeigen
						fmt.Printf("... und %d weitere Kontakte\n", len(contacts)-20)
						fmt.Println("💡 Verwenden Sie 'list all' für alle Kontakte")
						break
					}

					name := contact.PushName
					if name == "" {
						name = contact.BusinessName
					}
					if name == "" {
						name = "Unbekannt"
					}

					fmt.Printf("📞 %s (%s)\n", name, jid.User)
					count++
				}
			}
			fmt.Println()

		case "help", "h":
			fmt.Println("\n📋 Verfügbare Befehle:")
			fmt.Println("  📤 send <nummer/gruppe> <nachricht> - Nachricht senden")
			fmt.Println("     👤 Einzelperson: send 491234567890 Hallo")
			fmt.Println("     🏘️ Gruppe: send 120363XX...XX@g.us Hallo Gruppe")
			fmt.Println("  📋 list                            - Top 20 Kontakte anzeigen")
			fmt.Println("  📋 list all                        - Alle Kontakte anzeigen")
			fmt.Println("  🏘️ list groups                     - Gruppen anzeigen")
			fmt.Println("  🚪 quit                            - Programm beenden")
			fmt.Println("  ❓ help                            - Diese Hilfe anzeigen")

		default:
			fmt.Printf("❌ Unbekannter Befehl: %s (Verwenden Sie 'help' für Hilfe)\n", command)
		}
	}
}
