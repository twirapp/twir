export { getProfile, logout } from '@/services/auth/api.js';

export { createUserDashboard, selectedDashboardStore } from '@/services/auth/dashboard.js';
export { useTwitchAuth, useUserProfile } from '@/services/auth/hooks.js';
export {
  redirectToDashboard,
  redirectToLanding,
  redirectToLogin,
} from '@/services/auth/locationHelpers.js';
