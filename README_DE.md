# WhatsApp Console Client

Eine leistungsstarke Konsolenanwendung zum Senden von WhatsApp-Nachrichten über das Terminal - sowohl interaktiv als auch per CLI für Skripte.

## ✨ Features

- 🚀 **CLI-Einzeiler** für Skript-Integration  
- 💬 **Interaktiver Chat-Modus** für manuelle Nutzung
- 👤 **Einzelpersonen-Nachrichten** senden
- 🏘️ **Gruppen-Nachrichten** senden
- 📋 **Kontakte & Gruppen auflisten**
- 🔐 **WhatsApp Web QR-Code Login**
- 📱 **Empfangene Nachrichten anzeigen**
- 💾 **Session-Persistierung** in SQLite
- ⚡ **Schnell & Effizient**

## 📋 Voraussetzungen

- Go 1.21 oder höher
- SQLite3 (für Session-Speicherung)

## 🚀 Installation

```bash
cd whatsapp-console
make deps    # Abhängigkeiten installieren
make build   # Anwendung kompilieren
```

### **Cross-Platform Kompilierung**

Go ermöglicht es, für verschiedene Betriebssysteme und Architekturen zu kompilieren:

```bash
# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o whatsapp-console-linux-amd64 .

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o whatsapp-console-windows-amd64.exe .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o whatsapp-console-darwin-amd64 .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o whatsapp-console-darwin-arm64 .

# Linux (ARM64) - für Raspberry Pi, etc.
GOOS=linux GOARCH=arm64 go build -o whatsapp-console-linux-arm64 .
```

### **Alle Plattformen auf einmal erstellen**

```bash
# Skript erstellen, um für alle Plattformen zu kompilieren
mkdir -p dist

# Für alle wichtigen Plattformen kompilieren
for GOOS in linux windows darwin; do
    for GOARCH in amd64 arm64; do
        # ARM64 für Windows überspringen (nicht häufig verwendet)
        if [ "$GOOS" = "windows" ] && [ "$GOARCH" = "arm64" ]; then
            continue
        fi
        
        EXT=""
        if [ "$GOOS" = "windows" ]; then
            EXT=".exe"
        fi
        
        echo "Kompiliere für $GOOS/$GOARCH..."
        env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "dist/whatsapp-console-$GOOS-$GOARCH$EXT" .
    done
done
```

### **Docker-basierte Cross-Compilation**

Für konsistente Builds in verschiedenen Umgebungen:

```bash
# Dockerfile für das Kompilieren erstellen
cat > Dockerfile.build << 'EOF'
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Für mehrere Plattformen kompilieren
RUN mkdir -p /dist && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-linux-amd64 . && \
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /dist/whatsapp-console-linux-arm64 . && \
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-windows-amd64.exe . && \
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-darwin-amd64 . && \
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o /dist/whatsapp-console-darwin-arm64 .

FROM scratch
COPY --from=builder /dist /dist
EOF

# Mit Docker kompilieren
docker build -f Dockerfile.build -t whatsapp-console-builder .
docker run --rm -v $(pwd)/dist:/output whatsapp-console-builder sh -c "cp -r /dist/* /output/"
```

### **Plattform-spezifische Hinweise**

**🐧 Linux:**
- Native Kompilierung auf den meisten Distributionen
- ARM64-Builds funktionieren auf Raspberry Pi und ähnlichen Geräten

**🪟 Windows:**
- `.exe` Erweiterung wird automatisch hinzugefügt
- Cross-Compilation von Linux/macOS funktioniert perfekt

**🍎 macOS:**
- Intel (amd64) und Apple Silicon (arm64) unterstützt
- Code-Signierung kann für die Verteilung erforderlich sein

**🔧 Embedded Systems:**
```bash
# Für Raspberry Pi
GOOS=linux GOARCH=arm GOARM=7 go build -o whatsapp-console-raspberry .

# Für andere ARM-Geräte
GOOS=linux GOARCH=arm64 go build -o whatsapp-console-arm64 .
```

## 🎯 Verwendung

### **CLI-Modus (für Skripte)**

```bash
# Hilfe anzeigen
./whatsapp-console help

# Einzelperson
./whatsapp-console send 491234567890 'Hallo Welt!'

# Gruppe  
./whatsapp-console send 120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us 'Hallo Gruppe!'

# Mit Leerzeichen (in Anführungszeichen)
./whatsapp-console send 491234567890 'Nachricht mit mehreren Wörtern'
```

### **Interaktiver Modus**

```bash
./whatsapp-console
```

**Erste Ausführung:** QR-Code scannen erforderlich
1. Öffnen Sie WhatsApp auf Ihrem Handy
2. Gehen Sie zu Einstellungen → Verknüpfte Geräte  
3. Tippen Sie auf "Gerät verknüpfen"
4. Scannen Sie den QR-Code im Terminal

### **Interaktive Befehle**

- 📤 `send <nummer/gruppe> <nachricht>` - Nachricht senden
- 📋 `list` - Top 20 Kontakte anzeigen  
- 📋 `list all` - Alle Kontakte anzeigen
- 🏘️ `list groups` - Gruppen anzeigen
- ❓ `help` - Hilfe anzeigen
- 🚪 `quit` - Programm beenden

## 💡 Beispiele

### **CLI-Verwendung (Skripte)**

```bash
# Einzelperson benachrichtigen
./whatsapp-console send 491234567890 'Server ist wieder online!'

# Gruppe benachrichtigen  
./whatsapp-console send 120363420172074021@g.us 'Deployment erfolgreich!'

# Mit Variablen
NUMMER="491234567890"
NACHRICHT="Build #${BUILD_NUMBER} fertig"
./whatsapp-console send "$NUMMER" "$NACHRICHT"
```

### **Interaktive Verwendung**

```bash
# Anwendung starten
./whatsapp-console

# Im Chat-Modus:
whatsapp> send 491234567890 Hallo von der Konsole!
whatsapp> list
whatsapp> list groups  
whatsapp> help
whatsapp> quit
```

## 🔧 Entwicklung

### **Makefile Befehle**

```bash
make deps      # Abhängigkeiten installieren
make build     # Anwendung kompilieren
make run       # Anwendung ausführen
make clean     # Build-Artefakte löschen
make help      # Hilfe anzeigen
```

### **Projektstruktur**

```text
whatsapp-console/
├── main.go          # Hauptanwendung mit CLI & interaktivem Modus
├── go.mod           # Go Module Definition
├── go.sum           # Dependency Checksums
├── Makefile         # Build-Automatisierung
├── README.md        # Diese Dokumentation
└── whatsapp.db      # SQLite Session-Datenbank (wird erstellt)
```

## 🔧 Technische Details

- **WhatsApp API:** `go.mau.fi/whatsmeow` (Multidevice)
- **Database:** SQLite für Session-Persistierung
- **QR-Code:** Terminal-Anzeige mit `qrterminal`
- **Protokoll:** WhatsApp Web Protocol
- **Architektur:** Single-Binary mit dual-mode (CLI + Interactive)

## 🏘️ Gruppen-IDs finden

### **Methode 1: WhatsApp Web**
1. Öffnen Sie [web.whatsapp.com](https://web.whatsapp.com)
2. Klicken Sie auf eine Gruppe
3. Die URL enthält die ID: `chat/120363XXXXXXXXXX-XXXXXXXXXXXXXXXX`

### **Methode 2: Entwicklertools**  
1. WhatsApp Web → F12 → Console
2. Gruppe öffnen → ID aus der URL kopieren

### **Methode 3: Invite-Link**
- Gruppen-Einladungslinks enthalten die ID
- Format: `120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us`

## 🚨 Fehlerbehandlung

### **Häufige Probleme:**

**❌ "Erstes Login erforderlich"**
```bash
# Lösung: Interaktiv starten und QR-Code scannen
./whatsapp-console
# Dann CLI verwenden
```

**❌ "Timeout beim Verbinden"**
- Internetverbindung prüfen
- WhatsApp Web Status überprüfen
- Neuer Login: `rm whatsapp.db`

**❌ "Ungültige Nummer"**  
- Deutsche Nummern: `491234567890` (ohne +)
- Internationale: `41791234567` (ohne +)
- Gruppen: Vollständige ID mit `@g.us`

## 📱 Integration in Skripte

### **Bash-Beispiel**

```bash
#!/bin/bash
WHATSAPP="./whatsapp-console"
ADMIN="491234567890"
GROUP="120363420172074021@g.us"

# Server-Monitoring
if ! ping -c 1 google.com > /dev/null; then
    $WHATSAPP send $ADMIN "🚨 Server offline!"
fi

# Deployment-Benachrichtigung  
$WHATSAPP send $GROUP "✅ Release v1.2.3 deployed"
```

### **Cron-Job Beispiel**

```bash
# Täglich um 9:00 Uhr Status senden
0 9 * * * /path/to/whatsapp-console send 491234567890 "☀️ Guten Morgen! Server läuft."
```

## 🔐 Sicherheit

- ✅ **End-to-End Verschlüsselung** durch WhatsApp
- ✅ **Lokale Session-Speicherung** (verschlüsselt)
- ✅ **Keine Passwörter/Tokens** erforderlich  
- ✅ **Offizielle WhatsApp API** verwendet
- ⚠️ **Datei-Berechtigung** für `whatsapp.db` beachten

## 🆘 Support & Troubleshooting

### **Bei Problemen prüfen:**

1. **Go Version:** `go version` (mindestens 1.21)
2. **Abhängigkeiten:** `make deps`
3. **Kompilierung:** `make build` 
4. **Berechtigung:** `chmod +x whatsapp-console`
5. **Session:** Bei Problemen `rm whatsapp.db` und neu anmelden

### **Debug-Modus:**
```bash
# Mehr Logging für Fehleranalyse
export WHATSAPP_DEBUG=1
./whatsapp-console
```
## 🎯 Roadmap & Erweiterte Features

### **Geplante Features:**
- [ ] **Medien senden** (Bilder, Dokumente)
- [ ] **Gruppen verwalten** (erstellen, Mitglieder hinzufügen)
- [ ] **Nachrichtenverlauf** exportieren
- [ ] **Webhook-Integration** für eingehende Nachrichten
- [ ] **Multi-Account Support**
- [ ] **Config-Datei** für Einstellungen

### **Erweiterte Nutzung:**
- 🔄 **CI/CD Integration** für Deployment-Benachrichtigungen
- 📊 **Monitoring-Alerts** per WhatsApp
- 🤖 **Chatbot-Framework** aufbauen
- 📈 **Business-Integration** für Kundenservice

## 🤝 Beitragen

Contributions sind willkommen! 

```bash
# Repository forken und klonen
git clone https://github.com/your-username/whatsapp-console
cd whatsapp-console

# Feature-Branch erstellen  
git checkout -b feature/neue-funktion

# Änderungen committen
git commit -m "Neue Funktion hinzugefügt"

# Pull Request erstellen
```

## 📜 Lizenz

MIT License - Siehe LICENSE Datei für Details.

---

**🚀 Happy Messaging!** Gebaut mit ❤️ für die WhatsApp-Community.