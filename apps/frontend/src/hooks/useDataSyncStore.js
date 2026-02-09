import { createContext, useContext, useState } from 'react';

/**
 * Context para tracking de cambios pendientes entre vistas
 * Cuando se hacen cambios inline en Expenses, se marca como "dirty"
 * Dashboard verifica esto al montar y recarga datos si es necesario
 */
const DataSyncContext = createContext();

export const DataSyncProvider = ({ children }) => {
    const [expensesDirty, setExpensesDirty] = useState(false);
    const [incomesDirty, setIncomesDirty] = useState(false);

    const value = {
        // Flags de datos modificados pendientes de sincronización
        expensesDirty,
        incomesDirty,

        // Marcar expenses como modificados
        markExpensesDirty: () => setExpensesDirty(true),

        // Marcar incomes como modificados
        markIncomesDirty: () => setIncomesDirty(true),

        // Limpiar flag de expenses (después de sincronizar)
        clearExpensesDirty: () => setExpensesDirty(false),

        // Limpiar flag de incomes (después de sincronizar)
        clearIncomesDirty: () => setIncomesDirty(false),

        // Limpiar todos los flags
        clearAllDirty: () => {
            setExpensesDirty(false);
            setIncomesDirty(false);
        },
    };

    return (
        <DataSyncContext.Provider value={value}>
            {children}
        </DataSyncContext.Provider>
    );
};

export const useDataSync = () => {
    const context = useContext(DataSyncContext);
    if (!context) {
        throw new Error('useDataSync must be used within DataSyncProvider');
    }
    return context;
};

export default useDataSync;
