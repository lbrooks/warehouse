import { InventoryState } from './reducers/inventoryReducer'

interface AppStore {
  inventory: InventoryState
}

export const getInventoryState = (store: AppStore) => store.inventory;

export const getInventoryList = (store: AppStore) => getInventoryState(store) ? getInventoryState(store).items : [];
