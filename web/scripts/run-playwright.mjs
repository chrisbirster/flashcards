import net from 'node:net';
import { spawn } from 'node:child_process';

const PREFERRED_PORT = Number(process.env.PLAYWRIGHT_PORT || 5000);
const MAX_PORT_ATTEMPTS = 200;

function canConnect(host, port) {
  return new Promise((resolve) => {
    const socket = net.connect({ host, port });
    const done = (result) => {
      socket.removeAllListeners();
      socket.destroy();
      resolve(result);
    };

    socket.setTimeout(200);
    socket.on('connect', () => done(true));
    socket.on('timeout', () => done(false));
    socket.on('error', () => done(false));
  });
}

function canBind(host, port) {
  return new Promise((resolve) => {
    const server = net.createServer();
    server.unref();
    server.on('error', () => resolve(false));
    server.listen({ host, port }, () => {
      server.close(() => resolve(true));
    });
  });
}

async function isPortFree(port) {
  if (await canConnect('127.0.0.1', port)) return false;
  if (await canConnect('::1', port)) return false;
  if (!(await canBind('127.0.0.1', port))) return false;
  return canBind('::1', port);
}

async function findOpenPort(startPort) {
  for (let port = startPort; port < startPort + MAX_PORT_ATTEMPTS; port += 1) {
    if (await isPortFree(port)) return port;
  }
  throw new Error(`Unable to find an open port from ${startPort} to ${startPort + MAX_PORT_ATTEMPTS - 1}`);
}

const port = await findOpenPort(PREFERRED_PORT);
process.env.PLAYWRIGHT_PORT = String(port);
console.log(`[playwright] using frontend port ${port}`);

const extraArgs = process.argv.slice(2);
const command = process.platform === 'win32' ? 'npx.cmd' : 'npx';
const child = spawn(command, ['playwright', 'test', ...extraArgs], {
  stdio: 'inherit',
  env: process.env,
});

child.on('error', (error) => {
  console.error(`[playwright] failed to start test runner: ${error.message}`);
  process.exit(1);
});

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});
