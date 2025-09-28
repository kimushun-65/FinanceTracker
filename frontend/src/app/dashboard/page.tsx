import { DashboardContainer } from '../../page-components/dashboard';
import { ProtectedRoute } from '../../shared/ui/components/ProtectedRoute';

export default function Dashboard() {
  return (
    <ProtectedRoute>
      <DashboardContainer />
    </ProtectedRoute>
  );
}
