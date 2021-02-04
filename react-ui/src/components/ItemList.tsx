import {Items, Item} from '../models/items';

interface ItemListProps {
  items?: Items
}
const ItemList = (props: ItemListProps) => {
  const incr = (item: Item) => {
    alert(`Increment: ${JSON.stringify(item)}`);
	}

	const decr = (item: Item) => {
		alert(`Decrement: ${JSON.stringify(item)}`);
	}

  const { items } = props;
  if (!items || items.length === 0) return <p>No items, sorry</p>;
  return (
    <table>
      <thead>
        <tr>
          <th>Quantity</th>
          <th>Item</th>
          <th>Brand</th>
          <th>Barcode</th>
          <th>Add</th>
          {/* <th>Minus</th> */}
        </tr>
      </thead>
      <tbody>
        {items.map((i, idx) => {
          return (
            <tr key={'item-'+idx}>
              <td>{i.quantity}</td>
              <td>{i.name}</td>
              <td>{i.brand}</td>
              <td>{i.barcode}</td>
              <td><button type="button" onClick={() => incr(i)}>+</button></td>
              <td><button type="button" onClick={() => decr(i)}>-</button></td>
            </tr>
          );
        })}
      </tbody>
    </table>

  );
};

export default ItemList;