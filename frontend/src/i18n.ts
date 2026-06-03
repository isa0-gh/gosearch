import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      appName: "gosearch",
      search: "Search",
      placeholder: "Search anything...",
      tabs: {
        web: "Web",
        software: "Software",
        academic: "Academic",
        vuln: "Vulnerabilities",
        apps: "Apps",
        ml: "LLM/ML",
        games: "Games",
      },
      engine: "Engine",
      source: "Source",
      pages: "Pages",
      noResults: "No results found.",
      loading: "Searching...",
      error: "Something went wrong. Try again.",
      results: "results",
      nextPage: "Next page",
      fields: {
        downloads: "Downloads",
        size: "Size",
        stars: "Stars",
        language: "Language",
        updated: "Updated",
        category: "Category",
        pulls: "Pulls",
        tags: "Tags",
      },
      lang: "TR",
    },
  },
  tr: {
    translation: {
      appName: "gosearch",
      search: "Ara",
      placeholder: "Bir şeyler ara...",
      tabs: {
        web: "Web",
        software: "Yazılım",
        academic: "Akademik",
        vuln: "Güvenlik Açıkları",
        apps: "Uygulamalar",
        ml: "LLM/ML",
        games: "Oyunlar",
      },
      engine: "Motor",
      source: "Kaynak",
      pages: "Sayfa",
      noResults: "Sonuç bulunamadı.",
      loading: "Aranıyor...",
      error: "Bir hata oluştu. Tekrar deneyin.",
      results: "sonuç",
      nextPage: "Sonraki sayfa",
      fields: {
        downloads: "İndirme",
        size: "Boyut",
        stars: "Yıldız",
        language: "Dil",
        updated: "Güncellendi",
        category: "Kategori",
        pulls: "Çekme",
        tags: "Etiket",
      },
      lang: "EN",
    },
  },
};

i18n.use(initReactI18next).init({
  resources,
  lng: "en",
  fallbackLng: "en",
  interpolation: { escapeValue: false },
});

export default i18n;
