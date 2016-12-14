<<<if .flash.error>>>
<section class="flash error">
	<h3><<<.flash.error>>></h3>
	</section>
<<<end>>>
<<<if .Errors>>>
<section class="flash error">
	<ul class='list-unstyled'>
		<<<range $rec := .Errors>>>
		<li><<<$rec>>></li>
		<<<end>>>
	</ul>
</section>
<<<end>>>