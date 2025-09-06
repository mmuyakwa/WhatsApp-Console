# WhatsApp Console Client

A powerful console application for sending WhatsApp messages via terminal - both interactive and CLI modes for scripts.

## ✨ Features

- 🚀 **CLI one-liners** for script integration  
- 💬 **Interactive chat mode** for manual usage
- 👤 **Individual messages** to contacts
- 🏘️ **Group messages** support
- 📋 **List contacts & groups**
- 🔐 **WhatsApp Web QR-Code login**
- 📱 **Display received messages**
- 💾 **Session persistence** in SQLite
- ⚡ **Fast & Efficient**

## 📋 Prerequisites

- Go 1.21 or higher
- SQLite3 (for session storage)

## 🚀 Installation

```bash
cd whatsapp-console
make deps    # Install dependencies
make build   # Compile application
```

### **Cross-Platform Compilation**

Go allows you to compile for different operating systems and architectures:

```bash
# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o whatsapp-console-linux-amd64 ./cmd/main.go

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o whatsapp-console-windows-amd64.exe ./cmd/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o whatsapp-console-darwin-amd64 ./cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o whatsapp-console-darwin-arm64 ./cmd/main.go

# Linux (ARM64) - for Raspberry Pi, etc.
GOOS=linux GOARCH=arm64 go build -o whatsapp-console-linux-arm64 ./cmd/main.go
```

### **Build All Platforms at Once**

```bash
# Create a script to build for all platforms
mkdir -p dist

# Build for all major platforms
for GOOS in linux windows darwin; do
    for GOARCH in amd64 arm64; do
        # Skip ARM64 for Windows (not commonly used)
        if [ "$GOOS" = "windows" ] && [ "$GOARCH" = "arm64" ]; then
            continue
        fi
        
        EXT=""
        if [ "$GOOS" = "windows" ]; then
            EXT=".exe"
        fi
        
        echo "Building for $GOOS/$GOARCH..."
        env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "dist/whatsapp-console-$GOOS-$GOARCH$EXT" ./cmd/main.go
    done
done
```

### **Docker-based Cross-Compilation**

For consistent builds across different environments:

```bash
# Create a Dockerfile for building
cat > Dockerfile.build << 'EOF'
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build for multiple platforms
RUN mkdir -p /dist && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-linux-amd64 ./cmd/main.go && \
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /dist/whatsapp-console-linux-arm64 ./cmd/main.go && \
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-windows-amd64.exe ./cmd/main.go && \
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o /dist/whatsapp-console-darwin-amd64 ./cmd/main.go && \
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o /dist/whatsapp-console-darwin-arm64 ./cmd/main.go

FROM scratch
COPY --from=builder /dist /dist
EOF

# Build using Docker
docker build -f Dockerfile.build -t whatsapp-console-builder .
docker run --rm -v $(pwd)/dist:/output whatsapp-console-builder sh -c "cp -r /dist/* /output/"
```

### **Platform-Specific Notes**

**🐧 Linux:**
- Native compilation on most distributions
- ARM64 builds work on Raspberry Pi and similar devices

**🪟 Windows:**
- `.exe` extension is automatically added
- Cross-compilation from Linux/macOS works perfectly

**🍎 macOS:**
- Intel (amd64) and Apple Silicon (arm64) supported
- Code signing may be required for distribution

**🔧 Embedded Systems:**
```bash
# For Raspberry Pi
GOOS=linux GOARCH=arm GOARM=7 go build -o whatsapp-console-raspberry ./cmd/main.go

# For other ARM devices
GOOS=linux GOARCH=arm64 go build -o whatsapp-console-arm64 ./cmd/main.go
```

## 🎯 Usage

### **CLI Mode (for scripts)**

```bash
# Show help
./whatsapp-console help

# Individual contact
./whatsapp-console send 491234567890 'Hello World!'

# Group message  
./whatsapp-console send 120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us 'Hello Group!'

# With spaces (use quotes)
./whatsapp-console send 491234567890 'Message with multiple words'
```

### **Interactive Mode**

```bash
./whatsapp-console
```

**First run:** QR-Code scanning required
1. Open WhatsApp on your phone
2. Go to Settings → Linked Devices  
3. Tap "Link a Device"
4. Scan the QR-Code in the terminal

### **Interactive Commands**

- 📤 `send <number/group> <message>` - Send message
- 📋 `list` - Show top 20 contacts  
- 📋 `list all` - Show all contacts
- 🏘️ `list groups` - Show groups
- ❓ `help` - Show help
- 🚪 `quit` - Exit program

## 💡 Examples

### **CLI Usage (Scripts)**

```bash
# Notify individual contact
./whatsapp-console send 491234567890 'Server is back online!'

# Notify group  
./whatsapp-console send 120363420172074021@g.us 'Deployment successful!'

# With variables
NUMBER="491234567890"
MESSAGE="Build #${BUILD_NUMBER} completed"
./whatsapp-console send "$NUMBER" "$MESSAGE"
```

### **Interactive Usage**

```bash
# Start application
./whatsapp-console

# In chat mode:
whatsapp> send 491234567890 Hello from console!
whatsapp> list
whatsapp> list groups  
whatsapp> help
whatsapp> quit
```

## 🔧 Development

### **Makefile Commands**

```bash
make deps      # Install dependencies
make build     # Compile application
make run       # Run application
make clean     # Clean build artifacts
make help      # Show help
```

### **Project Structure**

```text
whatsapp-console/
├── main.go          # Main application with CLI & interactive mode
├── go.mod           # Go Module definition
├── go.sum           # Dependency checksums
├── Makefile         # Build automation
├── README.md        # This documentation
└── whatsapp.db      # SQLite session database (auto-created)
```

## 🔧 Technical Details

- **WhatsApp API:** `go.mau.fi/whatsmeow` (Multidevice)
- **Database:** SQLite for session persistence
- **QR-Code:** Terminal display with `qrterminal`
- **Protocol:** WhatsApp Web Protocol
- **Architecture:** Single-Binary with dual-mode (CLI + Interactive)

## 🏘️ Finding Group IDs

### **Method 1: WhatsApp Web**
1. Open [web.whatsapp.com](https://web.whatsapp.com)
2. Click on a group
3. The URL contains the ID: `chat/120363XXXXXXXXXX-XXXXXXXXXXXXXXXX`

### **Method 2: Developer Tools**  
1. WhatsApp Web → F12 → Console
2. Open group → copy ID from URL

### **Method 3: Invite Link**
- Group invitation links contain the ID
- Format: `120363XXXXXXXXXX-XXXXXXXXXXXXXXXX@g.us`

## 🚨 Error Handling

### **Common Issues:**

**❌ "First login required"**
```bash
# Solution: Start interactively and scan QR-Code
./whatsapp-console
# Then use CLI
```

**❌ "Connection timeout"**
- Check internet connection
- Verify WhatsApp Web status
- Fresh login: `rm whatsapp.db`

**❌ "Invalid number"**  
- German numbers: `491234567890` (without +)
- International: `41791234567` (without +)
- Groups: Full ID with `@g.us`

## 📱 Script Integration

### **Bash Example**

```bash
#!/bin/bash
WHATSAPP="./whatsapp-console"
ADMIN="491234567890"
GROUP="120363420172074021@g.us"

# Server monitoring
if ! ping -c 1 google.com > /dev/null; then
    $WHATSAPP send $ADMIN "🚨 Server offline!"
fi

# Deployment notification  
$WHATSAPP send $GROUP "✅ Release v1.2.3 deployed"
```

### **Cron Job Example**

```bash
# Daily at 9:00 AM send status
0 9 * * * /path/to/whatsapp-console send 491234567890 "☀️ Good morning! Server running."
```

## 🔐 Security

- ✅ **End-to-End encryption** via WhatsApp
- ✅ **Local session storage** (encrypted)
- ✅ **No passwords/tokens** required  
- ✅ **Official WhatsApp API** used
- ⚠️ **File permissions** for `whatsapp.db` recommended

## 🆘 Support & Troubleshooting

### **Check these when having issues:**

1. **Go Version:** `go version` (minimum 1.21)
2. **Dependencies:** `make deps`
3. **Compilation:** `make build` 
4. **Permissions:** `chmod +x whatsapp-console`
5. **Session:** For issues try `rm whatsapp.db` and re-login

### **Debug Mode:**
```bash
# More logging for error analysis
export WHATSAPP_DEBUG=1
./whatsapp-console
```

## 🎯 Roadmap & Advanced Features

### **Planned Features:**
- [ ] **Media sending** (images, documents)
- [ ] **Group management** (create, add members)
- [ ] **Message history** export
- [ ] **Webhook integration** for incoming messages
- [ ] **Multi-account support**
- [ ] **Config file** for settings

### **Advanced Usage:**
- 🔄 **CI/CD Integration** for deployment notifications
- 📊 **Monitoring alerts** via WhatsApp
- 🤖 **Chatbot framework** development
- 📈 **Business integration** for customer service

## 🤝 Contributing

Contributions are welcome! 

```bash
# Fork and clone repository
git clone https://github.com/your-username/whatsapp-console
cd whatsapp-console

# Create feature branch  
git checkout -b feature/new-feature

# Commit changes
git commit -m "Add new feature"

# Create pull request
```

## 📜 License

MIT License - See LICENSE file for details.

---

**🚀 Happy Messaging!** Built with ❤️ for the WhatsApp community.