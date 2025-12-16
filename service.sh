#!/bin/bash
#
# macOS Sensor Exporter Service Management Script
# Similar to 'brew services' but works without Homebrew
#

set -e

SERVICE_NAME="com.xykong.macos-sensor-exporter"
PLIST_FILE="${SERVICE_NAME}.plist"
LAUNCH_AGENTS_DIR="${HOME}/Library/LaunchAgents"
PLIST_PATH="${LAUNCH_AGENTS_DIR}/${PLIST_FILE}"
LOG_DIR="/usr/local/var/log"
LOG_FILE="${LOG_DIR}/macos-sensor-exporter.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_usage() {
    cat <<EOF
Usage: $0 <command>

Commands:
  install     Install the service (launchd plist)
  uninstall   Uninstall the service
  start       Start the service
  stop        Stop the service
  restart     Restart the service
  status      Show service status
  logs        Show service logs (tail -f)
  help        Show this help message

Examples:
  $0 install   # Install and start the service
  $0 start     # Start the service
  $0 status    # Check if service is running
  $0 logs      # View live logs

EOF
}

ensure_log_dir() {
    if [ ! -d "$LOG_DIR" ]; then
        echo "Creating log directory: $LOG_DIR"
        sudo mkdir -p "$LOG_DIR"
        sudo chmod 755 "$LOG_DIR"
    fi
}

install_service() {
    echo "Installing ${SERVICE_NAME} service..."
    
    # Create LaunchAgents directory if it doesn't exist
    mkdir -p "$LAUNCH_AGENTS_DIR"
    
    # Ensure log directory exists
    ensure_log_dir
    
    # Download or copy plist file
    if [ -f "$PLIST_FILE" ]; then
        cp "$PLIST_FILE" "$PLIST_PATH"
        echo -e "${GREEN}✓${NC} Copied plist file to $PLIST_PATH"
    else
        # Download from GitHub
        echo "Downloading plist file..."
        curl -fsSL "https://raw.githubusercontent.com/xykong/macos-sensor-exporter/master/${PLIST_FILE}" -o "$PLIST_PATH"
        echo -e "${GREEN}✓${NC} Downloaded plist file to $PLIST_PATH"
    fi
    
    # Load the service
    launchctl load "$PLIST_PATH" 2>/dev/null || true
    launchctl start "$SERVICE_NAME" 2>/dev/null || true
    
    echo -e "${GREEN}✓${NC} Service installed and started"
    echo ""
    echo "Service management commands:"
    echo "  $0 start    - Start the service"
    echo "  $0 stop     - Stop the service"
    echo "  $0 restart  - Restart the service"
    echo "  $0 status   - Check service status"
    echo "  $0 logs     - View service logs"
}

uninstall_service() {
    echo "Uninstalling ${SERVICE_NAME} service..."
    
    # Stop and unload the service
    launchctl stop "$SERVICE_NAME" 2>/dev/null || true
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    
    # Remove plist file
    if [ -f "$PLIST_PATH" ]; then
        rm "$PLIST_PATH"
        echo -e "${GREEN}✓${NC} Removed plist file"
    fi
    
    echo -e "${GREEN}✓${NC} Service uninstalled"
}

start_service() {
    if ! is_installed; then
        echo -e "${RED}✗${NC} Service not installed. Run: $0 install"
        exit 1
    fi
    
    echo "Starting ${SERVICE_NAME}..."
    launchctl load "$PLIST_PATH" 2>/dev/null || true
    launchctl start "$SERVICE_NAME"
    sleep 1
    
    if is_running; then
        echo -e "${GREEN}✓${NC} Service started"
    else
        echo -e "${RED}✗${NC} Failed to start service"
        exit 1
    fi
}

stop_service() {
    if ! is_installed; then
        echo -e "${RED}✗${NC} Service not installed"
        exit 1
    fi
    
    echo "Stopping ${SERVICE_NAME}..."
    launchctl stop "$SERVICE_NAME" 2>/dev/null || true
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    
    echo -e "${GREEN}✓${NC} Service stopped"
}

restart_service() {
    echo "Restarting ${SERVICE_NAME}..."
    stop_service
    sleep 1
    start_service
}

is_installed() {
    [ -f "$PLIST_PATH" ]
}

is_running() {
    launchctl list | grep -q "$SERVICE_NAME"
}

show_status() {
    echo "Service: ${SERVICE_NAME}"
    echo ""
    
    if is_installed; then
        echo -e "Installed: ${GREEN}✓${NC} Yes"
        echo "Plist: $PLIST_PATH"
    else
        echo -e "Installed: ${RED}✗${NC} No"
        echo "Run '$0 install' to install the service"
        return
    fi
    
    echo ""
    if is_running; then
        echo -e "Status: ${GREEN}✓${NC} Running"
        
        # Try to get PID
        PID=$(launchctl list | grep "$SERVICE_NAME" | awk '{print $1}')
        if [ -n "$PID" ] && [ "$PID" != "-" ]; then
            echo "PID: $PID"
        fi
    else
        echo -e "Status: ${RED}✗${NC} Not running"
    fi
    
    echo ""
    echo "Log file: $LOG_FILE"
    if [ -f "$LOG_FILE" ]; then
        echo "Log size: $(du -h "$LOG_FILE" | cut -f1)"
        echo ""
        echo "Recent logs:"
        tail -n 5 "$LOG_FILE" 2>/dev/null || echo "  (no logs yet)"
    fi
}

show_logs() {
    if [ ! -f "$LOG_FILE" ]; then
        echo -e "${YELLOW}⚠${NC} Log file not found: $LOG_FILE"
        exit 1
    fi
    
    echo "Showing logs from: $LOG_FILE"
    echo "Press Ctrl+C to exit"
    echo ""
    tail -f "$LOG_FILE"
}

# Main command dispatch
case "${1:-}" in
    install)
        install_service
        ;;
    uninstall)
        uninstall_service
        ;;
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    help|--help|-h|"")
        print_usage
        ;;
    *)
        echo -e "${RED}Error:${NC} Unknown command: $1"
        echo ""
        print_usage
        exit 1
        ;;
esac
