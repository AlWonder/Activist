<!DOCTYPE html>
</div> <!--The end of the container from the header-->
</main>
<footer class='main-footer'>
	<div class="container-fluid">
		<div class="row">
			<p>Разработано для ОСО БГПУ</p>
		</div>
	</div> 
</footer>

<script src='/js/jquery-3.1.1.min.js'></script>
<script src='/js/angular.min.js'></script>
<script src='/js/bootstrap.min.js'></script>
<script src='/js/app.js'></script>
<<< if .Js >>>
	<<< range $js := .Js >>>
	<script src='/js/<<< $js >>>.js'></script>
	<<< end >>>
<<< end >>>
</body>
</html>