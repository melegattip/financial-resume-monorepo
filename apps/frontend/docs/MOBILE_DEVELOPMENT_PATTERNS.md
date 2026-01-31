# 🎨 Financial Resume Engine - Patrones de Desarrollo Mobile

## 🧩 Componentes Reutilizables

### 1. **ValidatedInput** - Input con Validación
```typescript
interface ValidatedInputProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  validation: (value: string) => { isValid: boolean; error: string | null };
  type?: 'text' | 'number' | 'email' | 'password';
  placeholder?: string;
  required?: boolean;
}

// Uso en formularios
<ValidatedInput
  label="Monto"
  value={formData.amount}
  onChange={(value) => setFormData({...formData, amount: value})}
  validation={validateAmount}
  type="number"
  placeholder="0.00"
  required
/>
```

### 2. **ConfirmationModal** - Modal de Confirmación
```typescript
interface ConfirmationModalProps {
  isVisible: boolean;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
  variant?: 'danger' | 'warning' | 'info';
}

// Uso para eliminaciones
<ConfirmationModal
  isVisible={showDeleteModal}
  title="Eliminar Gasto"
  message="¿Estás seguro de que quieres eliminar este gasto? Esta acción no se puede deshacer."
  confirmText="Eliminar"
  onConfirm={confirmDelete}
  onCancel={() => setShowDeleteModal(false)}
  loading={deleteLoading}
  variant="danger"
/>
```

### 3. **LoadingState** - Estados de Carga
```typescript
interface LoadingStateProps {
  loading: boolean;
  error: string | null;
  onRetry?: () => void;
  children: React.ReactNode;
}

// Uso universal para páginas
<LoadingState loading={loading} error={error} onRetry={loadData}>
  {/* Contenido de la página */}
</LoadingState>
```

### 4. **AmountDisplay** - Mostrar Montos
```typescript
interface AmountDisplayProps {
  amount: number;
  hidden?: boolean;
  variant?: 'income' | 'expense' | 'neutral';
  size?: 'sm' | 'md' | 'lg';
  showCurrency?: boolean;
}

// Uso en dashboard y listas
<AmountDisplay
  amount={expense.amount}
  hidden={balancesHidden}
  variant="expense"
  size="md"
  showCurrency
/>
```

### 5. **CategoryBadge** - Badge de Categoría
```typescript
interface CategoryBadgeProps {
  category: Category | null;
  variant?: 'default' | 'outline';
  size?: 'sm' | 'md';
}

<CategoryBadge 
  category={expense.category} 
  variant="outline" 
  size="sm" 
/>
```

---

## 📱 Patrones de Navegación Mobile

### 1. **Bottom Tab Navigation**
```typescript
const TabNavigator = () => (
  <Tab.Navigator
    screenOptions={{
      tabBarStyle: {
        backgroundColor: theme.colors.background,
        borderTopColor: theme.colors.border,
      },
      tabBarActiveTintColor: theme.colors.primary,
      tabBarInactiveTintColor: theme.colors.textSecondary,
    }}
  >
    <Tab.Screen 
      name="Dashboard" 
      component={DashboardScreen}
      options={{
        tabBarIcon: ({ color, size }) => <HomeIcon color={color} size={size} />,
        title: 'Inicio'
      }}
    />
    <Tab.Screen 
      name="Expenses" 
      component={ExpensesScreen}
      options={{
        tabBarIcon: ({ color, size }) => <MinusCircleIcon color={color} size={size} />,
        title: 'Gastos'
      }}
    />
    <Tab.Screen 
      name="Incomes" 
      component={IncomesScreen}
      options={{
        tabBarIcon: ({ color, size }) => <PlusCircleIcon color={color} size={size} />,
        title: 'Ingresos'
      }}
    />
    <Tab.Screen 
      name="Analytics" 
      component={ReportsScreen}
      options={{
        tabBarIcon: ({ color, size }) => <ChartBarIcon color={color} size={size} />,
        title: 'Análisis'
      }}
    />
    <Tab.Screen 
      name="Profile" 
      component={AchievementsScreen}
      options={{
        tabBarIcon: ({ color, size }) => <UserIcon color={color} size={size} />,
        title: 'Perfil'
      }}
    />
  </Tab.Navigator>
);
```

### 2. **Stack Navigation**
```typescript
const MainStackNavigator = () => (
  <Stack.Navigator
    screenOptions={{
      headerStyle: {
        backgroundColor: theme.colors.background,
      },
      headerTintColor: theme.colors.text,
      headerTitleStyle: {
        fontWeight: 'bold',
      },
    }}
  >
    <Stack.Screen 
      name="Main" 
      component={TabNavigator} 
      options={{ headerShown: false }}
    />
    <Stack.Screen 
      name="ExpenseDetail" 
      component={ExpenseDetailScreen}
      options={({ route }) => ({ 
        title: route.params?.expense?.description || 'Detalle del Gasto' 
      })}
    />
    <Stack.Screen 
      name="CreateExpense" 
      component={CreateExpenseScreen}
      options={{ 
        title: 'Nuevo Gasto',
        presentation: 'modal'
      }}
    />
  </Stack.Navigator>
);
```

---

## 🎯 Flujos de Trabajo Principales

### 1. **Flujo de Autenticación**
```typescript
const AuthFlow = {
  // 1. Check initial auth state
  initial: async () => {
    const token = await SecureStore.getItemAsync('auth_token');
    const user = await SecureStore.getItemAsync('auth_user');
    
    if (token && user) {
      const isValid = await validateToken(token);
      if (isValid) {
        return { authenticated: true, user: JSON.parse(user) };
      }
    }
    
    return { authenticated: false, user: null };
  },

  // 2. Login process
  login: async (credentials) => {
    const response = await api.post('/auth/login', credentials);
    const { token, user, expires_at } = response.data.data;
    
    // Store securely
    await SecureStore.setItemAsync('auth_token', token);
    await SecureStore.setItemAsync('auth_user', JSON.stringify(user));
    await SecureStore.setItemAsync('auth_expires_at', expires_at.toString());
    
    return { user, token };
  },

  // 3. Logout process
  logout: async () => {
    try {
      await api.post('/auth/logout');
    } catch (error) {
      // Continue with logout even if API call fails
    }
    
    // Clear stored data
    await SecureStore.deleteItemAsync('auth_token');
    await SecureStore.deleteItemAsync('auth_user');
    await SecureStore.deleteItemAsync('auth_expires_at');
    
    // Clear app cache
    await AsyncStorage.clear();
  }
};
```

### 2. **Flujo de Sincronización de Datos**
```typescript
const DataSyncFlow = {
  // 1. Load data with cache support
  loadWithCache: async (endpoint, cacheKey, maxAge = 300000) => {
    // Try cache first
    const cached = await AsyncStorage.getItem(`cache_${cacheKey}`);
    if (cached) {
      const { data, timestamp } = JSON.parse(cached);
      if (Date.now() - timestamp < maxAge) {
        return { data, source: 'cache' };
      }
    }
    
    // Fetch from API
    try {
      const response = await api.get(endpoint);
      const data = response.data;
      
      // Cache the result
      await AsyncStorage.setItem(`cache_${cacheKey}`, JSON.stringify({
        data,
        timestamp: Date.now()
      }));
      
      return { data, source: 'api' };
    } catch (error) {
      // Return cached data if available, even if expired
      if (cached) {
        const { data } = JSON.parse(cached);
        return { data, source: 'cache_fallback' };
      }
      throw error;
    }
  },

  // 2. Invalidate cache after mutations
  invalidateCache: async (pattern) => {
    const keys = await AsyncStorage.getAllKeys();
    const cacheKeys = keys.filter(key => 
      key.startsWith('cache_') && key.includes(pattern)
    );
    
    await AsyncStorage.multiRemove(cacheKeys);
  },

  // 3. Sync data in background
  backgroundSync: async () => {
    // Only sync if user is authenticated
    const token = await SecureStore.getItemAsync('auth_token');
    if (!token) return;
    
    try {
      // Critical data first
      await Promise.all([
        DataSyncFlow.loadWithCache('/dashboard', 'dashboard'),
        DataSyncFlow.loadWithCache('/gamification/profile', 'gamification'),
      ]);
      
      // Secondary data
      await Promise.all([
        DataSyncFlow.loadWithCache('/categories', 'categories'),
        DataSyncFlow.loadWithCache('/budgets/dashboard', 'budgets'),
      ]);
    } catch (error) {
      console.log('Background sync failed:', error);
    }
  }
};
```

### 3. **Flujo de Operaciones CRUD**
```typescript
const CRUDFlow = {
  // Generic CRUD operations with optimistic updates
  create: async (endpoint, data, cacheKey) => {
    // Optimistic update
    const tempId = `temp_${Date.now()}`;
    const optimisticItem = { ...data, id: tempId, __optimistic: true };
    
    // Update local state immediately
    updateLocalState(cacheKey, (items) => [optimisticItem, ...items]);
    
    try {
      const response = await api.post(endpoint, data);
      const newItem = response.data.data;
      
      // Replace optimistic item with real data
      updateLocalState(cacheKey, (items) =>
        items.map(item => item.id === tempId ? newItem : item)
      );
      
      // Invalidate related cache
      await DataSyncFlow.invalidateCache(cacheKey);
      
      return newItem;
    } catch (error) {
      // Remove optimistic item on error
      updateLocalState(cacheKey, (items) =>
        items.filter(item => item.id !== tempId)
      );
      throw error;
    }
  },

  update: async (endpoint, data, cacheKey, itemId) => {
    // Store original for rollback
    const originalItem = getLocalItem(cacheKey, itemId);
    
    // Optimistic update
    updateLocalState(cacheKey, (items) =>
      items.map(item => item.id === itemId ? { ...item, ...data } : item)
    );
    
    try {
      const response = await api.patch(`${endpoint}/${itemId}`, data);
      const updatedItem = response.data.data;
      
      // Update with server response
      updateLocalState(cacheKey, (items) =>
        items.map(item => item.id === itemId ? updatedItem : item)
      );
      
      return updatedItem;
    } catch (error) {
      // Rollback on error
      updateLocalState(cacheKey, (items) =>
        items.map(item => item.id === itemId ? originalItem : item)
      );
      throw error;
    }
  },

  delete: async (endpoint, itemId, cacheKey) => {
    // Store original for rollback
    const originalItem = getLocalItem(cacheKey, itemId);
    
    // Optimistic delete
    updateLocalState(cacheKey, (items) =>
      items.filter(item => item.id !== itemId)
    );
    
    try {
      await api.delete(`${endpoint}/${itemId}`);
      // Success - item already removed from local state
    } catch (error) {
      // Rollback on error
      updateLocalState(cacheKey, (items) => [originalItem, ...items]);
      throw error;
    }
  }
};
```

---

## 🎮 Gamificación Mobile

### 1. **Notificaciones de Gamificación**
```typescript
const GamificationNotifications = {
  showXPGained: (xp, action) => {
    // Show toast notification
    Toast.show({
      type: 'success',
      text1: `+${xp} XP`,
      text2: action,
      visibilityTime: 3000,
    });
    
    // Optional: Haptic feedback
    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
  },

  showLevelUp: (newLevel, levelName) => {
    // Show modal or full-screen animation
    showModal({
      title: '¡Nivel Subido!',
      message: `Has alcanzado el nivel ${newLevel}: ${levelName}`,
      type: 'celebration',
    });
    
    // Strong haptic feedback
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Heavy);
  },

  showAchievementUnlocked: (achievement) => {
    Toast.show({
      type: 'success',
      text1: '🏆 Logro Desbloqueado',
      text2: achievement.name,
      visibilityTime: 4000,
    });
  }
};
```

### 2. **Feature Gates Mobile**
```typescript
const FeatureGates = {
  checkAccess: async (featureKey) => {
    try {
      const response = await api.get(`/gamification/features/${featureKey}/access`);
      return response.data;
    } catch (error) {
      // Fallback to local logic
      const userProfile = await getUserProfile();
      const feature = FEATURE_GATES[featureKey];
      
      return {
        unlocked: userProfile.current_level >= feature.requiredLevel,
        requiredLevel: feature.requiredLevel,
        userLevel: userProfile.current_level,
        xpNeeded: Math.max(0, feature.xpThreshold - userProfile.total_xp)
      };
    }
  },

  showLockedFeature: (featureAccess) => {
    Alert.alert(
      'Función Bloqueada',
      `Necesitas alcanzar el nivel ${featureAccess.requiredLevel} para desbloquear esta función. Te faltan ${featureAccess.xpNeeded} XP.`,
      [
        { text: 'Entendido', style: 'cancel' },
        { text: 'Ver Progreso', onPress: () => navigate('Achievements') }
      ]
    );
  }
};
```

---

## 📊 Charts y Visualizaciones

### 1. **Chart Components para Mobile**
```typescript
// Pie Chart para categorías
const CategoryPieChart = ({ data, size = 200 }) => (
  <PieChart
    data={data}
    width={size}
    height={size}
    chartConfig={{
      color: (opacity = 1) => `rgba(26, 255, 146, ${opacity})`,
    }}
    accessor="amount"
    backgroundColor="transparent"
    paddingLeft="15"
    absolute
  />
);

// Bar Chart para trends
const TrendBarChart = ({ data, width, height = 220 }) => (
  <BarChart
    data={data}
    width={width}
    height={height}
    chartConfig={{
      backgroundColor: '#ffffff',
      backgroundGradientFrom: '#ffffff',
      backgroundGradientTo: '#ffffff',
      decimalPlaces: 0,
      color: (opacity = 1) => `rgba(0, 158, 227, ${opacity})`,
    }}
    verticalLabelRotation={30}
  />
);
```

### 2. **Dashboard Widgets Mobile**
```typescript
const DashboardWidget = ({ title, value, icon, color, onPress }) => (
  <TouchableOpacity 
    style={[styles.widget, { borderLeftColor: color }]}
    onPress={onPress}
  >
    <View style={styles.widgetHeader}>
      <Text style={styles.widgetTitle}>{title}</Text>
      {icon}
    </View>
    <Text style={[styles.widgetValue, { color }]}>{value}</Text>
  </TouchableOpacity>
);

// Usage
<DashboardWidget
  title="Balance Total"
  value={formatCurrency(balance)}
  icon={<DollarSignIcon size={24} color={Colors.primary} />}
  color={balance >= 0 ? Colors.success : Colors.error}
  onPress={() => navigate('Reports')}
/>
```

---

## 🔒 Seguridad Mobile

### 1. **Secure Storage**
```typescript
import * as SecureStore from 'expo-secure-store';
import * as LocalAuthentication from 'expo-local-authentication';

const SecureStorage = {
  // Store sensitive data
  setSecure: async (key, value) => {
    await SecureStore.setItemAsync(key, value, {
      requireAuthentication: true,
      authenticationPrompt: 'Autentícate para acceder a tus datos financieros'
    });
  },

  // Retrieve sensitive data
  getSecure: async (key) => {
    return await SecureStore.getItemAsync(key, {
      requireAuthentication: true,
      authenticationPrompt: 'Autentícate para acceder a tus datos financieros'
    });
  },

  // Check biometric availability
  checkBiometrics: async () => {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    const isEnrolled = await LocalAuthentication.isEnrolledAsync();
    const supportedTypes = await LocalAuthentication.supportedAuthenticationTypesAsync();
    
    return {
      available: hasHardware && isEnrolled,
      types: supportedTypes
    };
  }
};
```

### 2. **App State Security**
```typescript
const AppStateSecurity = {
  // Hide content when app goes to background
  setupAppStateListener: () => {
    AppState.addEventListener('change', (nextAppState) => {
      if (nextAppState === 'background') {
        // Hide sensitive content
        setShowSensitiveData(false);
        
        // Optional: Show privacy screen
        showPrivacyScreen();
      } else if (nextAppState === 'active') {
        // App came back to foreground
        hidePrivacyScreen();
        
        // Optional: Require re-authentication if app was in background too long
        checkSessionValidity();
      }
    });
  },

  // Session timeout management
  setupSessionTimeout: () => {
    let timeoutId;
    
    const resetTimeout = () => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => {
        // Auto-logout after inactivity
        logout();
      }, 15 * 60 * 1000); // 15 minutes
    };
    
    // Reset timeout on user interaction
    const panResponder = PanResponder.create({
      onStartShouldSetPanResponder: () => {
        resetTimeout();
        return false;
      }
    });
    
    return panResponder;
  }
};
```

---

## 🎨 Tema y Estilos

### 1. **Theme System**
```typescript
const lightTheme = {
  colors: {
    primary: '#007AFF',
    background: '#FFFFFF',
    surface: '#F8F9FA',
    text: '#000000',
    textSecondary: '#6C757D',
    border: '#DEE2E6',
    error: '#DC3545',
    success: '#28A745',
    warning: '#FFC107',
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
  },
  borderRadius: {
    sm: 4,
    md: 8,
    lg: 12,
    xl: 16,
  }
};

const darkTheme = {
  ...lightTheme,
  colors: {
    ...lightTheme.colors,
    primary: '#0A84FF',
    background: '#000000',
    surface: '#1C1C1E',
    text: '#FFFFFF',
    textSecondary: '#8E8E93',
    border: '#38383A',
  }
};
```

### 2. **Responsive Utilities**
```typescript
const ResponsiveUtils = {
  // Get device info
  getDeviceInfo: () => {
    const { width, height } = Dimensions.get('window');
    return {
      width,
      height,
      isSmall: width < 375,
      isMedium: width >= 375 && width < 414,
      isLarge: width >= 414,
      isTablet: width >= 768,
    };
  },

  // Scale font sizes
  scaleFontSize: (size) => {
    const { width } = Dimensions.get('window');
    const scale = width / 375; // Base width iPhone 8
    return Math.round(size * scale);
  },

  // Dynamic spacing
  getDynamicSpacing: (baseSpacing) => {
    const { isSmall, isTablet } = ResponsiveUtils.getDeviceInfo();
    
    if (isTablet) return baseSpacing * 1.5;
    if (isSmall) return baseSpacing * 0.8;
    return baseSpacing;
  }
};
```

---

**Este documento complementa ENDPOINTS_MOBILE_REFERENCE.md proporcionando los patrones de implementación específicos para desarrollo mobile nativo.** 