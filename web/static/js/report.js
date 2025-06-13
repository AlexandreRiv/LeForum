function toggleReportForm(id) {
    const form = document.getElementById('report-form-' + id);
    if (form.classList.contains('hidden')) {
        // Fermer tous les autres formulaires
        document.querySelectorAll('.report-form').forEach(f => {
            f.classList.add('hidden');
        });
        // Ouvrir celui-ci
        form.classList.remove('hidden');
    } else {
        form.classList.add('hidden');
    }
}