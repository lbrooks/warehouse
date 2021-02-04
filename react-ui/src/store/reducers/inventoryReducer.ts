import { InventoryActionTypes, LOAD_INVENTORY, LOAD_INVENTORY_SUCCESS, INCREMENT_INVENTORY, DECREMENT_INVENTORY } from "../types/inventoryTypes";
import { Item } from '../../models/items';
import { retrieveInventory } from "../../backend/inv-server";
import { loadInventorySuccess } from "../actions";
import store from '../store'

export interface InventoryState {
  loading: boolean
  items: Item[]
}

const initialState = {
  loading: false,
  items: []
} as InventoryState ;

export function inventoryReducer(state = initialState, action: InventoryActionTypes) {
  switch (action.type) {
    case LOAD_INVENTORY: {
      console.log(`into load`);
      // Make an API call to fetch todos from the server
      retrieveInventory('').then((items) => {
        // Dispatch an action with the todos we received
         store.dispatch(loadInventorySuccess(items))
      })

      return {
        loading: true,
        items: []
      };
    }
    case LOAD_INVENTORY_SUCCESS: {
      return {
        loading: false,
        items: action.payload.items
      };
    }
    case INCREMENT_INVENTORY: {
      return {
        ...state,
        items: state.items.map((i) => {
          if (action.payload.item.barcode === i.barcode && action.payload.item.brand === i.brand && action.payload.item.name === i.name) {
            return {...i, quantity: i.quantity + 1};
          }
          return i;
        })
      };
    }
    case DECREMENT_INVENTORY: {
      return {
        ...state,
        items: state.items.map((i) => {
          if (action.payload.item.barcode === i.barcode && action.payload.item.brand === i.brand && action.payload.item.name === i.name) {
            if (i.quantity > 0) {
              return {...i, quantity: i.quantity - 1};
            }
          }
          return i;
        })
      };
    }
    default:
      return state;
  }
}
