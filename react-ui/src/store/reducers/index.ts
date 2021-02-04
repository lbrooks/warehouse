import { combineReducers } from 'redux';
import { inventoryReducer } from './inventoryReducer';

export const rootReducer = combineReducers({
    inventory: inventoryReducer
});

export type RootState = ReturnType<typeof rootReducer>