import WebSocket from 'ws'
import * as readline from 'readline'
import * as crypto from 'crypto'
import type { ClientMessage } from '../shared/protocol.js'
import {
	isBroadcastMessage,
	isConnectedMessage,
	isUserJoinedMessage,
	isUserLeftMessage,
	isErrorMessage,
	isAckMessage,
	createSendMessage,
	createDisconnectMessage,
} from '../shared/protocol.js'
import { MessageRouter } from '../shared/message-router.js'

const url = process.argv[2]
const userId = process.argv[3] || crypto.randomUUID()

if (!url) {
	console.error('Usage: npx ts-node client.ts <websocket-url> [user-id]')
	console.error('Example: npx ts-node client.ts wss://xxx.apigw.yandexcloud.net/ws')
	console.error('Example: npx ts-node client.ts wss://xxx.apigw.yandexcloud.net/ws my-user-123')
	process.exit(1)
}

// Append user_id to URL query params
const wsUrl = new URL(url)
wsUrl.searchParams.set('user_id', userId)

console.log(`User ID: ${userId}`)
console.log(`Connecting to ${wsUrl.toString()}...`)

const ws = new WebSocket(wsUrl.toString())

function send(msg: ClientMessage): void {
	ws.send(JSON.stringify(msg))
}

// ANSI color codes
const green = '\x1b[32m'
const reset = '\x1b[0m'

// Track pending messages for acknowledgment
let pendingMessage: string | null = null

// Print a message while preserving the pending line
function printMessage(text: string, stream: NodeJS.WriteStream = process.stdout): void {
	stream.write(text + '\n')
}

const rl = readline.createInterface({
		input: process.stdin,
		output: process.stdout,
		prompt: '> ',
	})

// Configure message router with handlers for each message type
const router = new MessageRouter()
	.on(isBroadcastMessage, (msg) => {
		printMessage(`< [${msg.from}]: ${msg.message}`)
	})
	.on(isConnectedMessage, (msg) => {
		printMessage(`* Connected as ${msg.userId}`, process.stderr)
		rl.prompt();
	})
	.on(isUserJoinedMessage, (msg) => {
		printMessage(`* ${msg.userId} joined`, process.stderr)
	})
	.on(isUserLeftMessage, (msg) => {
		printMessage(`* ${msg.userId} left`, process.stderr)
	})
	.on(isErrorMessage, (msg) => {
		process.stderr.write(`< Error: ${msg.message}\n`)
	})
	.on(isAckMessage, (_msg) => {
		// When ack received, replace "> " with green "✓ "
		if (pendingMessage !== null) {
			// Carriage return, clear line, print checkmark, then newline
			process.stdout.write(`\r\x1b[K${green}✓${reset} ${pendingMessage}\n`)
			pendingMessage = null
			rl.prompt()
		}
	})
	.onUnknown((msg) => {
		printMessage(`< ${JSON.stringify(msg)}`)
	})
	.onParseError((rawData) => {
		printMessage(`< Received: ${rawData}`)
	})

ws.on('open', () => {
	console.log('Connected!')
	console.log('Type messages and press Enter to send. Press Ctrl+C to exit.')
	console.log('---')

	rl.on('line', (line) => {
		const message = line.trim()
		if (message) {
			const sendMsg = createSendMessage(message)

			pendingMessage = message
			send(sendMsg)
			// Clear the line with user input, write "? message" with carriage return
			process.stdout.write(`\x1b[A\r\x1b[K? ${message}`)
		}
	})

	rl.on('close', () => {
		console.log('\nSending disconnect...')
		const disconnectMsg = createDisconnectMessage()
		send(disconnectMsg)
		ws.close()
		process.exit(0)
	})
})

ws.on('message', (data) => router.handle(data))

ws.on('close', (code, reason) => {
	console.log(`Disconnected (code: ${code}, reason: ${reason.toString() || 'none'})`)
	process.exit(0)
})

ws.on('error', (error) => {
	console.error('WebSocket error:', error.message)
	process.exit(1)
})
