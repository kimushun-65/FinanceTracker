import { AssetsContainer } from '../../page-components/assets';
import { ProtectedRoute } from '../../shared/ui/components/ProtectedRoute';

export default function AssetsPage() {
  return (
    <ProtectedRoute>
      <AssetsContainer />
    </ProtectedRoute>
  );
}
