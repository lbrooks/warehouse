import {Item, Items} from '../../models/items';

export const LOAD_INVENTORY = "inventory/load";
export const LOAD_INVENTORY_SUCCESS = "inventory/load/success";
export const INCREMENT_INVENTORY = "inventory/increment";
export const DECREMENT_INVENTORY = "inventory/decrement";

interface LoadInventoryAction {
  type: typeof LOAD_INVENTORY
}

interface LoadInventorySuccessAction {
  type: typeof LOAD_INVENTORY_SUCCESS
  payload: {
    items: Items
  }
}

interface IncrementInventoryAction {
  type: typeof INCREMENT_INVENTORY
  payload: {
    item: Item
  }
}

interface DecrementInventoryAction {
  type: typeof DECREMENT_INVENTORY
  payload: {
    item: Item
  }
}

export type InventoryActionTypes = LoadInventoryAction | LoadInventorySuccessAction | IncrementInventoryAction | DecrementInventoryAction
