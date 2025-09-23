import { DashboardContainer } from '../../page-components/dashboard';
import { ProtectedRoute } from '../../shared/components/ProtectedRoute';

export default function Dashboard() {
  return (
    <ProtectedRoute>
      <DashboardContainer />
    </ProtectedRoute>
  );
}
