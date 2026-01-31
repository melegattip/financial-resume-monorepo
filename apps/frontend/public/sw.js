const CACHE_NAME = 'finance-app-v1.0.0';
const urlsToCache = [
  '/',
  '/static/js/bundle.js',
  '/static/css/main.css',
  '/manifest.json',
  '/icon-192x192.png',
  '/icon-512x512.png'
];

// Instalación del Service Worker
self.addEventListener('install', (event) => {
  console.log('Service Worker: Instalando...');
  
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('Service Worker: Cache abierto');
        return cache.addAll(urlsToCache);
      })
      .then(() => {
        console.log('Service Worker: Instalación completada');
        return self.skipWaiting();
      })
  );
});

// Activación del Service Worker
self.addEventListener('activate', (event) => {
  console.log('Service Worker: Activando...');
  
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            console.log('Service Worker: Eliminando cache viejo:', cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    }).then(() => {
      console.log('Service Worker: Activación completada');
      return self.clients.claim();
    })
  );
});

// Intercepción de peticiones (Fetch)
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        // Devolver desde cache si existe
        if (response) {
          return response;
        }

        // Si no está en cache, hacer petición a la red
        return fetch(event.request).then((response) => {
          // Verificar si es una respuesta válida
          if (!response || response.status !== 200 || response.type !== 'basic') {
            return response;
          }

          // Clonar la respuesta para cachearla
          const responseToCache = response.clone();

          caches.open(CACHE_NAME)
            .then((cache) => {
              cache.put(event.request, responseToCache);
            });

          return response;
        });
      })
      .catch(() => {
        // Fallback para páginas offline
        if (event.request.destination === 'document') {
          return caches.match('/offline.html');
        }
      })
  );
});

// Manejo de notificaciones push
self.addEventListener('push', (event) => {
  console.log('Service Worker: Push recibido', event);

  const options = {
    body: 'Tienes nuevas actualizaciones financieras',
    icon: '/icon-192x192.png',
    badge: '/badge-72x72.png',
    vibrate: [200, 100, 200],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    },
    actions: [
      {
        action: 'explore',
        title: 'Ver Detalles',
        icon: '/icons/explore.png'
      },
      {
        action: 'close',
        title: 'Cerrar',
        icon: '/icons/close.png'
      }
    ]
  };

  event.waitUntil(
    self.registration.showNotification('FinanceApp', options)
  );
});

// Manejo de clics en notificaciones
self.addEventListener('notificationclick', (event) => {
  console.log('Service Worker: Notification click', event);

  const notification = event.notification;
  const action = event.action;

  if (action === 'close') {
    notification.close();
  } else {
    // Abrir o enfocar la aplicación
    event.waitUntil(
      clients.matchAll().then((clientList) => {
        for (const client of clientList) {
          if (client.url === '/' && 'focus' in client) {
            return client.focus();
          }
        }
        
        if (clients.openWindow) {
          return clients.openWindow('/');
        }
      })
    );
  }

  notification.close();
});

// Sincronización en background
self.addEventListener('sync', (event) => {
  console.log('Service Worker: Background sync', event);

  if (event.tag === 'background-sync') {
    event.waitUntil(
      syncData()
    );
  }
});

// Función para sincronizar datos offline
async function syncData() {
  try {
    // Obtener datos pendientes del IndexedDB
    const pendingData = await getPendingData();
    
    for (const item of pendingData) {
      // Intentar enviar cada item pendiente
      await fetch('/api/sync', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(item)
      });
      
      // Eliminar del storage local si se envió exitosamente
      await removePendingData(item.id);
    }
    
    console.log('Service Worker: Sincronización completada');
  } catch (error) {
    console.error('Service Worker: Error en sincronización:', error);
  }
}

// Simulación de funciones de IndexedDB
async function getPendingData() {
  // Implementar lógica para obtener datos pendientes
  return [];
}

async function removePendingData(id) {
  // Implementar lógica para eliminar datos sincronizados
  console.log('Eliminando datos sincronizados:', id);
}

// Manejo de mensajes desde la aplicación
self.addEventListener('message', (event) => {
  console.log('Service Worker: Mensaje recibido', event.data);

  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }

  if (event.data && event.data.type === 'SCHEDULE_NOTIFICATION') {
    const { title, body, delay } = event.data;
    
    setTimeout(() => {
      self.registration.showNotification(title, {
        body,
        icon: '/icon-192x192.png',
        badge: '/badge-72x72.png'
      });
    }, delay);
  }
}); 