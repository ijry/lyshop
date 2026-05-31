/**
 * eapp mock data — re-exports admin mock for build compatibility.
 *
 * Why this file exists:
 * - eapp/utils/request.ts dynamically imports '../../admin/src/mock/index'
 * - Vite cannot bundle cross-directory dynamic imports correctly
 * - This file provides a local entry point that re-exports admin mock
 * - Allows Vite to include mock data in the build output
 */

export { matchMock } from '../../admin/src/mock/index'
