// components/LoadingSpinner.tsx
export default function LoadingSpinner() {
  return (
    <div class="flex justify-center items-center py-12">
      <div class="relative">
        <div class="w-12 h-12 rounded-full absolute border-4 border-gray-200"></div>
        <div class="w-12 h-12 rounded-full animate-spin absolute border-4 border-blue-500 border-t-transparent"></div>
      </div>
      <span class="ml-4 text-gray-600 font-medium">Loading books...</span>
    </div>
  );
}