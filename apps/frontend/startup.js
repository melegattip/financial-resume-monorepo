#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

// Función para detectar el ambiente
function detectEnvironment() {
  // 1. Si NODE_ENV está explícitamente configurado
  if (process.env.NODE_ENV === 'production') {
    return 'production';
  }
  
  if (process.env.NODE_ENV === 'development') {
    return 'development';
  }
  
  // 2. Si REACT_APP_ENVIRONMENT está configurado
  if (process.env.REACT_APP_ENVIRONMENT === 'production') {
    return 'production';
  }
  
  // 3. Detectar por estructura de archivos
  const hasSrcDir = fs.existsSync(path.join(__dirname, 'src'));
  const hasBuildDir = fs.existsSync(path.join(__dirname, 'build'));
  const hasNodeModules = fs.existsSync(path.join(__dirname, 'node_modules'));
  
  // Si tenemos build pero no src, probablemente es producción
  if (hasBuildDir && !hasSrcDir) {
    return 'production';
  }
  
  // Si tenemos src, probablemente es desarrollo
  if (hasSrcDir) {
    return 'development';
  }
  
  // 4. Detectar por hostname si estamos en Render
  const hostname = process.env.RENDER_SERVICE_NAME || '';
  if (hostname.includes('render') || process.env.RENDER === 'true') {
    return 'production';
  }
  
  // Fallback: si tenemos build, usar producción; sino desarrollo
  return hasBuildDir ? 'production' : 'development';
}

// Función para ejecutar el comando correcto
function startApp() {
  const environment = detectEnvironment();
  const port = process.env.PORT || 3000;
  
  console.log(`🔧 Startup Script: Detected environment: ${environment}`);
  console.log(`🚀 Starting on port: ${port}`);
  
  let command, args;
  
  if (environment === 'production') {
    // Producción: servir build estático
    console.log('📦 Starting production server with serve...');
    command = 'npx';
    args = ['serve', '-s', 'build', '-l', port.toString()];
  } else {
    // Desarrollo: servidor de React con hot reload
    console.log('🔧 Starting development server with react-scripts...');
    command = 'npx';
    args = ['react-scripts', 'start'];
    
    // En desarrollo, configurar puerto
    process.env.PORT = port.toString();
  }
  
  // Ejecutar el comando
  const child = spawn(command, args, {
    stdio: 'inherit',
    shell: true,
    env: process.env
  });
  
  child.on('error', (error) => {
    console.error(`❌ Error starting app: ${error.message}`);
    process.exit(1);
  });
  
  child.on('close', (code) => {
    console.log(`🏁 App process exited with code ${code}`);
    process.exit(code);
  });
  
  // Manejar señales para shutdown graceful
  process.on('SIGINT', () => {
    console.log('\n🛑 Received SIGINT, shutting down gracefully...');
    child.kill('SIGINT');
  });
  
  process.on('SIGTERM', () => {
    console.log('\n🛑 Received SIGTERM, shutting down gracefully...');
    child.kill('SIGTERM');
  });
}

// Iniciar la aplicación
startApp(); 