function submit_playbook() {
	$.post(
		"/submit",
		{
			content: $('#playbook_content').val()
		},
		update_results
	).fail(handle_failure);
	$('#results').hide();
	$('#loading').show();

	return false;
}

function update_results(data) {
	$('#loading').hide();
	$('#results').html(data);
	$('#results').show();
}

function handle_failure(data) {
	update_results(data.responseText);
}

$(document).ready(function() {
	$('#playbook_form').submit(submit_playbook);
	$('#playbook_submit').click(submit_playbook);
});
