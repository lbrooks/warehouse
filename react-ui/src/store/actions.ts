import { InventoryActionTypes, LOAD_INVENTORY, LOAD_INVENTORY_SUCCESS, INCREMENT_INVENTORY, DECREMENT_INVENTORY} from "./types/inventoryTypes";
import { Item, Items } from '../models/items';

// TypeScript infers that this function is returning SendMessageAction
export function loadInventory(): InventoryActionTypes {
  return {type: LOAD_INVENTORY};
}

// TypeScript infers that this function is returning DeleteMessageAction
export function loadInventorySuccess(items: Items): InventoryActionTypes {
  return {
    type: LOAD_INVENTORY_SUCCESS,
    payload: {items}
  };
}

// TypeScript infers that this function is returning DeleteMessageAction
export function incrementItem(item: Item): InventoryActionTypes {
  return {
    type: INCREMENT_INVENTORY,
    payload: {item}
  };
}

// TypeScript infers that this function is returning DeleteMessageAction
export function decrementItem(item: Item): InventoryActionTypes {
  return {
    type: DECREMENT_INVENTORY,
    payload: {item}
  };
}
