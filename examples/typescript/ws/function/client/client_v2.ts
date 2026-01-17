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
	console.error('Usage: npx ts-node client_v2.ts <websocket-url> [user-id]')
	console.error('Example: npx ts-node client_v2.ts wss://xxx.apigw.yandexcloud.net/ws')
	console.error('Example: npx ts-node client_v2.ts wss://xxx.apigw.yandexcloud.net/ws my-user-123')
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
const blue = '\x1b[34m'
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
		// In v2, we receive our own messages back via the trigger
		// Display them differently to distinguish from other users
		if (msg.from === userId) {
			// Our own message coming back via broadcast trigger
			if (pendingMessage !== null && msg.message === pendingMessage) {
				// This is the echo of the message we just sent
				// Clear the pending "?" indicator and show green checkmark
				process.stdout.write(`\r\x1b[K${green}✓${reset} ${msg.message}\n`)
				pendingMessage = null
			} else {
				// Our message from another session/tab (same userId)
				printMessage(`${blue}< [You (other session)]:${reset} ${msg.message}`)
			}
		} else {
			// Message from another user
			printMessage(`< [${msg.from}]: ${msg.message}`)
		}
		rl.prompt()
	})
	.on(isConnectedMessage, (msg) => {
		printMessage(`* Connected as ${msg.userId}`, process.stderr)
		rl.prompt();
	})
	.on(isUserJoinedMessage, (msg) => {
		// Don't show "you joined" for ourselves
		if (msg.userId !== userId) {
			printMessage(`* ${msg.userId} joined`, process.stderr)
		}
	})
	.on(isUserLeftMessage, (msg) => {
		// Don't show "you left" for ourselves (shouldn't happen anyway)
		if (msg.userId !== userId) {
			printMessage(`* ${msg.userId} left`, process.stderr)
		}
	})
	.on(isErrorMessage, (msg) => {
		process.stderr.write(`< Error: ${msg.message}\n`)
	})
	.on(isAckMessage, (_msg) => {
		// ACK just confirms message was written to topic
		// We'll show the green checkmark when we receive the broadcast echo
		// Keep the pending message for now
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
	console.log('Note: Messages are broadcast via Data Streams trigger in { messages: [...] } format')
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

	setInterval(() => {
		ws.ping()
	}, 10000);
})

ws.on('message', (data, isBinary) => router.handle(data, isBinary))

ws.on('close', (code, reason) => {
	console.log(`Disconnected (code: ${code}, reason: ${reason.toString() || 'none'})`)
	process.exit(0)
})

ws.on('error', (error) => {
	console.error('WebSocket error:', error.message)
	process.exit(1)
})
