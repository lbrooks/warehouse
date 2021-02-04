import {Items} from '../models/items';

export function retrieveInventory(barcode: string): Promise<Items> {
    return fetch(`/api/item?barcode=${barcode}`)
        .then((res) => res.json());
}
