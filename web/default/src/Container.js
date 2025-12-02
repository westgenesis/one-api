import { useLocation } from 'react-router-dom';
import { Container } from 'semantic-ui-react';

import App from './App';
const LayoutContainer = () => {
  const { pathname } = useLocation();

  const isMonitor = pathname === '/monitor';

  return (
    <Container className={isMonitor ? 'monitor-content' : 'main-content'}>
      <App />
    </Container>
  );
};

export default LayoutContainer;
