interface WithLoadingProps {
  loading: boolean
}

function WithLoading(Component: any) {
  return function WihLoadingComponent({ loading, ...props }: WithLoadingProps) {
    if (!loading) return <Component {...props} />;
    return (
      <p style={{ textAlign: 'center', fontSize: '30px' }}>
        Hold on, fetching data may take some time :)
      </p>
    );
  };
}
export default WithLoading;