import { useEffect } from 'react'
import { useSelector } from 'react-redux'
import './App.css';
import ItemList from './components/ItemList';
import withLoading from './components/WithLoading';
import { loadInventory } from './store/actions';
import { getInventoryState } from './store/selectors';
import store from './store/store';

function App() {
  const ListLoading = withLoading(ItemList);
  const invState = useSelector(getInventoryState)

  // Note: the empty deps array [] means
  // this useEffect will run once
  // similar to componentDidMount()
  useEffect(() => {
    store.dispatch(loadInventory())
  }, [])

  return (
    <div className="App">
      <ListLoading {...invState}></ListLoading>
    </div>
  );
}

export default App;
